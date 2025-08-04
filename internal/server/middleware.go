package server

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type Middleware struct {
	logger    *slog.Logger
	authToken string
}

func NewMiddleware(logger *slog.Logger, token string) *Middleware {
	return &Middleware{
		logger:    logger,
		authToken: token,
	}
}

func (m *Middleware) Register(router *echo.Echo) {
	router.Use(m.AccessLog())
}

func (m *Middleware) AccessLog() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			startTime := time.Now()
			requestID := uuid.New().String()
			c.Set("requestID", requestID)

			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				m.logger.Warn("Missing or malformed Authorization header",
					slog.String("RequestID", requestID),
					slog.String("IP", c.RealIP()),
					slog.String("URL", c.Request().URL.Path),
				)
				return echo.NewHTTPError(401, "unauthorized: missing or malformed Authorization header")
			}

			token := strings.TrimPrefix(authHeader, "Bearer ")
			if token != m.authToken {
				m.logger.Warn("Invalid token",
					slog.String("RequestID", requestID),
					slog.String("IP", c.RealIP()),
					slog.String("URL", c.Request().URL.Path),
					slog.String("ProvidedToken", token),
				)
				return echo.NewHTTPError(401, "unauthorized: invalid token")
			}

			m.logger.Info("Request started",
				slog.String("RequestID", requestID),
				slog.String("IP", c.RealIP()),
				slog.String("URL", c.Request().URL.Path),
				slog.String("Method", c.Request().Method),
			)

			err := next(c)
			responseTime := time.Since(startTime)

			if err != nil {
				m.logger.Error("Request Failed",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
					slog.String("Error", err.Error()),
				)
			} else {
				m.logger.Info("Request done",
					slog.String("RequestID", requestID),
					slog.String("Time spent", strconv.FormatInt(int64(responseTime), 10)),
				)
			}

			return err
		}
	}
}
