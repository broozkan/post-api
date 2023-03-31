package services

import (
	"context"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	RepositoryInterface interface {
		CreatePost(post *models.Post) error
		GetRankedPosts(offset, limit int, params map[string]string) ([]*models.Post, error)
		GetPromotedPosts() ([]*models.Post, error)
		GetTotalPostsCount() (int, error)
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
	post.ID = uuid.New().String()

	if err := s.repository.CreatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostsWithFilters(offset, limit int, params map[string]string) (*models.ListPostsResponse, error) {
	panic("implement me!")
}
