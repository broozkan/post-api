package handlers

import (
	"context"
	"net/http"
	"strconv"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

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
	h.logger.Debug("GetFeed request arrived")

	pageStr := c.Query("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		return c.Status(http.StatusBadRequest).JSON(&models.ErrorHandler{Message: "Invalid page number"})
	}

	offset := (page - 1) * h.config.ItemPerPage
	limit := h.config.ItemPerPage

	filterMap := parseQueryStringParams(c.Request().URI().QueryString())

	response, err := h.postService.GetPostsWithFilters(offset, limit, filterMap)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(&models.ErrorHandler{Message: "Failed to list posts"})
	}

	return c.Status(http.StatusOK).JSON(response)
}
