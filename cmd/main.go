package main

import (
	"log"
	"os"

	"broozkan/postapi/internal/config"

	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	appEnv := os.Getenv("APP_ENV")
	logger, _ := zap.NewProduction()
	defer func() {
		_ = logger.Sync()
	}()
	conf, err := config.New("../.config", appEnv)
	if err != nil {
		return err
	}
	logger.Info("config loaded", zap.Any("config", conf))
	return nil
}
