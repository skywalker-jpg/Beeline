package server

import (
	"TestBeeline/internal/models"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/labstack/echo"
)

func (s *Server) RegisterHandlers() {
	app := s.app

	apiGroup := app.Group("/api/v1")

	apiGroup.POST("/process", s.ProcessXML)

	app.GET("/*", s.NotFound)
	app.POST("/*", s.NotFound)
}

func (s *Server) ProcessXML(c echo.Context) error {
	requestID := c.Get("requestID").(string)

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		s.logger.Error("Failed to read request body",
			slog.String("RequestID", requestID),
			slog.String("Error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to read request body",
		})
	}

	var users models.Users
	if err := xml.Unmarshal(body, &users); err != nil {
		s.logger.Error("Failed to parse XML",
			slog.String("RequestID", requestID),
			slog.String("Error", err.Error()))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid XML format",
		})
	}

	s.logger.Info("XML parsed successfully",
		slog.String("RequestID", requestID),
		slog.Int("UsersCount", len(users.Users)))

	results := make(chan models.UserJSON, len(users.Users))
	errors := make(chan error, len(users.Users))
	var wg sync.WaitGroup

	for _, user := range users.Users {
		wg.Add(1)
		go func(u models.User) {
			defer wg.Done()
			s.processUser(u, results, requestID)
		}(user)
	}

	go func() {
		wg.Wait()
		close(results)
		close(errors)
	}()

	var processedUsers []models.UserJSON
	var processingErrors []error

	for {
		select {
		case user, ok := <-results:
			if !ok {
				results = nil
			} else {
				processedUsers = append(processedUsers, user)
			}
		case err, ok := <-errors:
			if !ok {
				errors = nil
			} else {
				processingErrors = append(processingErrors, err)
			}
		}

		if results == nil && errors == nil {
			break
		}
	}

	if len(processingErrors) > 0 {
		s.logger.Error("Errors during user processing",
			slog.String("RequestID", requestID),
			slog.Int("ErrorCount", len(processingErrors)))
		for _, err := range processingErrors {
			s.logger.Error("Processing error detail",
				slog.String("RequestID", requestID),
				slog.String("Error", err.Error()))
		}
	}

	// Отправка данных на удаленный сервер
	if len(processedUsers) > 0 {
		if err := s.sendToRemoteServer(processedUsers, requestID); err != nil {
			s.logger.Error("Failed to send data to remote server",
				slog.String("RequestID", requestID),
				slog.String("Error", err.Error()))
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to send data to remote server",
			})
		}
	}

	// Успешный ответ
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":        "Processing completed",
		"users_received": len(users.Users),
		"users_sent":     len(processedUsers),
		"errors":         len(processingErrors),
	})
}

// processUser обрабатывает одного пользователя
func (s *Server) processUser(user models.User, results chan<- models.UserJSON, requestID string) {
	s.logger.Debug("Processing user",
		slog.String("RequestID", requestID),
		slog.String("UserID", user.ID),
		slog.String("UserName", user.Name))

	// Преобразование возраста в возрастную группу
	ageGroup := s.getAgeGroup(user.Age)

	// Создание JSON объекта
	userJSON := models.UserJSON{
		ID:       user.ID,
		FullName: user.Name,
		Email:    user.Email,
		AgeGroup: ageGroup,
	}

	results <- userJSON
}

func (s *Server) getAgeGroup(age int) string {
	switch {
	case age < 25:
		return "до 25"
	case age >= 25 && age <= 35:
		return "от 25 до 35"
	default:
		return "старше 35"
	}
}

func (s *Server) sendToRemoteServer(users []models.UserJSON, requestID string) error {
	jsonData, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	s.logger.Info("Sending data to remote server",
		slog.String("RequestID", requestID),
		slog.Int("UsersCount", len(users)),
		slog.Int("DataSize", len(jsonData)))

	req, err := http.NewRequest("POST", s.serverURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", requestID)

	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Чтение ответа
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	s.logger.Info("Remote server response",
		slog.String("RequestID", requestID),
		slog.Int("StatusCode", resp.StatusCode),
		slog.String("Response", string(responseBody)))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("remote server returned status %d: %s", resp.StatusCode, string(responseBody))
	}

	s.logger.Info("Data sent successfully to remote server",
		slog.String("RequestID", requestID),
		slog.Int("UsersSent", len(users)))

	return nil
}

func (s *Server) NotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "Resource not found",
	})
}
