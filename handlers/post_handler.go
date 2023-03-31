package handlers

import (
	"context"
	"errors"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/asaskevich/govalidator"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type (
	PostServiceInterface interface {
		CreatePost(ctx context.Context, post *models.Post) (*models.Post, error)
		GetPostsWithFilters(offset, limit int, params map[string]string) (*models.ListPostsResponse, error)
	}

	PostHandler struct {
		logger      *zap.Logger
		postService PostServiceInterface
		config      *config.Config
		nsfwPosts   map[int]bool
	}
)

func NewPostHandler(logger *zap.Logger, conf *config.Config, postService PostServiceInterface) *PostHandler {
	return &PostHandler{
		logger:      logger,
		postService: postService,
		config:      conf,
		nsfwPosts:   nil,
	}
}

func (h *PostHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/posts", h.CreatePostHandler)
	app.Get("/posts/feed", h.GetFeedHandler)
}

func (h *PostHandler) CreatePostHandler(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := validatePost(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	post, err := h.postService.CreatePost(c.Context(), post)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create post",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

func (h *PostHandler) GetFeedHandler(c *fiber.Ctx) error {
	panic("implement me!")
}

func validatePost(post *models.Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}
	if post.Link != "" && !govalidator.IsURL(post.Link) {
		return errors.New("invalid link")
	}
	if post.Link != "" && post.Content != "" {
		return errors.New("a post cannot have both a link and content populated")
	}
	return nil
}
