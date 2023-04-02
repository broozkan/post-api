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

	conf, err := config.New("../.config", appEnv)
	if err != nil {
		return err
	}

	logger, err := initLogger(conf.LogLevel)
	if err != nil {
		return err
	}

	defer func() {
		_ = logger.Sync()
	}()

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

func initLogger(logLevel string) (*zap.Logger, error) {
	if logLevel == "" {
		logLevel = "debug"
	}
	var level zap.AtomicLevel
	if err := level.UnmarshalText([]byte(logLevel)); err != nil {
		return nil, err
	}

	conf := zap.Config{
		Level:            level,
		Development:      true,
		Encoding:         "console",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}

	return logger, nil
}
