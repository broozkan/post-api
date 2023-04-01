package server

import (
	"os"
	"os/signal"
	"syscall"

	"broozkan/postapi/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

type Handler interface {
	RegisterRoutes(app *fiber.App)
}

type Server struct {
	App    *fiber.App
	config config.Server
	logger *zap.Logger
}

func New(logger *zap.Logger, serverConfig config.Server, handlers []Handler) Server {
	app := fiber.New()

	server := Server{App: app, config: serverConfig, logger: logger}
	server.App.Use(cors.New())
	server.addRoutes()

	for _, handler := range handlers {
		handler.RegisterRoutes(server.App)
	}

	return server
}

func (s Server) Run() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		shutdownSignal := <-shutdownChan
		s.logger.Info("Received interrupt signal", zap.String("shutdownSignal", shutdownSignal.String()))
		if err := s.App.Shutdown(); err != nil {
			s.logger.Info("Failed to shutdown gracefully", zap.Error(err))
			return
		}
		s.logger.Info("application shutdown gracefully")
	}()
	err := s.App.Listen(s.config.Port)
	if err != nil {
		s.logger.Panic(err.Error())
	}
}

func (s Server) Stop() {
	err := s.App.Shutdown()
	if err != nil {
		s.logger.Info("Graceful shutdown failed")
	}
}

func (s Server) addRoutes() {
	s.App.Get("/health", healthCheck)
}

func healthCheck(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK)
	return nil
}
