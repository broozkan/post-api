package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/repository"
	"broozkan/postapi/internal/services"
	"broozkan/postapi/pkg/server"

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

	couchbaseRepo, err := repository.NewPostRepository(&conf.Couchbase)
	if err != nil {
		logger.Error("error while initializing couchbase", zap.Error(err))
		return err
	}

	postService := services.NewPostService(logger, conf, couchbaseRepo)

	postHandler := handlers.NewPostHandler(logger, conf, postService)

	serverHandlers := []server.Handler{
		postHandler,
	}

	s := server.New(logger, conf.Server, serverHandlers)

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	s.Run()

	return nil
}
