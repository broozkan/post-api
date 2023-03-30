package handlers

import (
	"context"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type (
	PostServiceInterface interface {
		CreatePost(ctx context.Context, post *models.Post) error
	}

	PostHandler struct {
		logger      *zap.Logger
		postService PostServiceInterface
		config      *config.Config
	}
)

func NewPostHandler(logger *zap.Logger, conf *config.Config, postService PostServiceInterface) *PostHandler {
	return &PostHandler{
		logger:      logger,
		postService: postService,
		config:      conf,
	}
}

func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/posts", h.CreatePostHandler)
}

func (h *PostHandler) CreatePostHandler(c *fiber.Ctx) error {
	panic("implement me!")
}
