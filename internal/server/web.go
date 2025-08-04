package server

import (
	"TestBeeline/internal/config"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"log/slog"
	"net/http"
)

type Server struct {
	app       *echo.Echo
	URL       string
	logger    *slog.Logger
	client    *http.Client
	serverURL string
}

func New(srvCfg config.Server, logger *slog.Logger) (*Server, error) {
	e := echo.New()
	server := Server{
		app:       e,
		URL:       srvCfg.URL,
		logger:    logger,
		client:    &http.Client{},
		serverURL: srvCfg.ServerURL,
	}
	e.HideBanner = true
	e.Logger.SetOutput(io.Discard)

	e.Use(middleware.Recover())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())

	m := NewMiddleware(logger, srvCfg.AuthToken)
	m.Register(e)
	server.RegisterHandlers()

	return &server, nil
}

func (s *Server) Serve() error {
	s.logger.Info("HTTP server started", slog.String("url", s.URL))

	return fmt.Errorf("server error: %w", s.app.Start(s.URL))
}
