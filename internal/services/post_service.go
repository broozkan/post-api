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
		GetPromotedPosts(count int) ([]*models.Post, error)
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
	posts, err := s.repository.GetRankedPosts(offset, limit, params)
	if err != nil {
		return nil, err
	}

	if len(posts) >= s.conf.MinPostLengthForAd && s.conf.AdsEnabled {
		var promotedPosts []*models.Post
		promotedPosts, err = s.repository.GetPromotedPosts(len(s.conf.AdsPositions))
		if err != nil {
			s.logger.Error("unable to get promoted posts", zap.Error(err)) // keep program running
		}

		if len(promotedPosts) != 0 {
			adIndices := PrepareIndices(posts, s.conf.AdsPositions)
			for _, idx := range adIndices {
				promotedPost := promotedPosts[generateRandomNumber(int64(len(promotedPosts)))]
				posts, err = AddPromotedPost(posts, promotedPost, idx)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	totalCount, err := s.repository.GetTotalPostsCount()
	if err != nil {
		return nil, err
	}

	totalPages := (totalCount + limit - 1) / limit
	response := &models.ListPostsResponse{
		Posts:      posts,
		Page:       offset/s.conf.ItemPerPage + 1,
		TotalPages: totalPages,
	}

	return response, nil
}
