package main

import (
	"TestBeeline/internal/config"
	"TestBeeline/internal/logger"
	"TestBeeline/internal/server"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	config, err := config.NewConfig("config.yaml")
	if err != nil {
		return fmt.Errorf("new config error: %w", err)
	}

	logger, err := logger.New(config.Logger)
	if err != nil {
		return fmt.Errorf("new logger error: %w", err)
	}

	server, err := server.New(config.Server, logger)
	if err != nil {
		return err
	}

	return server.Serve()
}
