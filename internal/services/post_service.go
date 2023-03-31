package services

import (
	"context"
	"crypto/rand"
	"math/big"

	"broozkan/postapi/handlers"
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/google/uuid"
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
	post.ID = uuid.New().String()

	if err := s.repository.CreatePost(post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetPostsWithFilters(offset, limit int, params map[string]string) (*models.ListPostsResponse, error) {
	panic("implement me!")
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		result[i] = letters[n.Int64()]
	}
	return string(result)
}
