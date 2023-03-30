package services

import (
	"context"
	"math/rand"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"go.uber.org/zap"
)

type (
	RepositoryInterface interface {
		CreatePost(post *models.Post) error
	}

	PostService struct {
		logger     *zap.Logger
		conf       *config.Config
		repository RepositoryInterface
	}
)

var _ handlers.PostServiceInterface = (*PostService)(nil)

func NewPostService(logger *zap.Logger, conf *config.Config, repository RepositoryInterface) *PostService {
	return &PostService{
		logger:     logger,
		conf:       conf,
		repository: repository,
	}
}

func (s *PostService) CreatePost(_ context.Context, post *models.Post) (*models.Post, error) {
	post.Author = s.conf.AuthorPrefix + randomString(s.conf.AuthorIDLength)

	err := s.repository.CreatePost(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
