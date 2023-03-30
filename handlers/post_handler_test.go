package handlers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/mocks"
	"broozkan/postapi/internal/models"

	"github.com/go-faker/faker/v4"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const testURL = "http://testdomain.com"

func TestPostHandler_CreatePostHandlerHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockService := mocks.NewMockPostServiceInterface(ctrl)

	app := fiber.New(fiber.Config{
		ErrorHandler: fiber.DefaultErrorHandler,
	})

	conf := &config.Config{}
	handler := handlers.NewPostHandler(zap.NewNop(), conf, mockService)
	handler.RegisterRoutes(app)

	t.Run("given valid posts request when service ok then it should return status ok", func(t *testing.T) {
		var post models.Post
		_ = faker.FakeData(&post)
		post.Link = testURL
		post.Content = "" // to prevent both content and link error
		mockService.EXPECT().CreatePost(gomock.Any(), gomock.Any()).Return(&post, nil)

		bodyBytes, _ := json.Marshal(&post)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Equal(t, fiber.StatusCreated, res.StatusCode)
	})

	t.Run("given invalid posts request when service ok then it should return status bad request", func(t *testing.T) {
		var post models.Post
		_ = faker.FakeData(&post)
		post.Author = ""

		bodyBytes, _ := json.Marshal(&post)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})

	t.Run("given valid posts request when service returns error then it should return internal server error", func(t *testing.T) {
		var post models.Post
		_ = faker.FakeData(&post)
		post.Link = testURL
		post.Author = "t2_11qnzrqv"
		post.Content = "" // to prevent both content and link error

		mockService.EXPECT().CreatePost(gomock.Any(), &post).Return(nil, errors.New("dummy error"))

		bodyBytes, _ := json.Marshal(&post)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Equal(t, fiber.StatusInternalServerError, res.StatusCode)
	})

	t.Run("given both link and content populated request when service ok then it should return status bad request", func(t *testing.T) {
		var post models.Post
		_ = faker.FakeData(&post)
		post.Link = testURL
		post.Content = "example content"

		bodyBytes, _ := json.Marshal(&post)
		req := httptest.NewRequest(http.MethodPost, "/posts", bytes.NewBuffer(bodyBytes))
		req.Header.Set("Content-Type", "application/json")
		res, err := app.Test(req)
		assert.Nil(t, err)
		res.Body.Close()

		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)
	})
}
