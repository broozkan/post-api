package main

import (
	"log"
	"os"

	"go.uber.org/zap"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	_ = os.Getenv("APP_ENV")
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	return nil
}
