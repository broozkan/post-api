package services

import (
	"context"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/models"

	"go.uber.org/zap"
)

type (
	RepositoryInterface interface {
		CreatePost(post *models.Post) error
	}

	PostService struct {
		logger     *zap.Logger
		repository RepositoryInterface
	}
)

var _ handlers.PostServiceInterface = (*PostService)(nil)

func NewPostService(logger *zap.Logger, repository RepositoryInterface) *PostService {
	return &PostService{
		logger:     logger,
		repository: repository,
	}
}

func (s *PostService) CreatePost(ctx context.Context, post *models.Post) error {
	panic("implement me!")
}
