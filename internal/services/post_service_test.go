package services_test

import (
	"context"
	"errors"
	"testing"

	"broozkan/postapi/internal/mocks"
	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/services"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostService_CreatePost(t *testing.T) {
	t.Run("given post request when repository is ok then it should return nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepositoryInterface(ctrl)

		var post models.Post
		_ = faker.FakeData(&post)

		mockRepository.EXPECT().CreatePost(gomock.Any()).Return(nil)

		postService := services.NewPostService(zap.NewNop(), mockRepository)
		err := postService.CreatePost(context.Background(), &post)
		assert.Nil(t, err)
	})

	t.Run("given post request when repository return error then it should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepositoryInterface(ctrl)

		var post models.Post
		_ = faker.FakeData(&post)

		mockRepository.EXPECT().CreatePost(gomock.Any()).Return(errors.New("dummy error"))

		postService := services.NewPostService(zap.NewNop(), mockRepository)
		err := postService.CreatePost(context.Background(), &post)
		assert.Nil(t, err)
	})
}
