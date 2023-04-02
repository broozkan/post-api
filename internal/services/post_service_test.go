package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/mocks"
	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/services"

	"github.com/go-faker/faker/v4"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestPostService_CreatePost(t *testing.T) {
	conf := &config.Config{}
	t.Run("given post request when repository is ok then it should return nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepositoryInterface(ctrl)

		var post models.Post
		_ = faker.FakeData(&post)

		mockRepository.EXPECT().CreatePost(gomock.Any()).Return(nil)

		postService := services.NewPostService(zap.NewNop(), conf, mockRepository)
		_, err := postService.CreatePost(context.Background(), &post)
		assert.Nil(t, err)
	})

	t.Run("given post request when repository return error then it should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mockRepository := mocks.NewMockRepositoryInterface(ctrl)

		var post models.Post
		_ = faker.FakeData(&post)

		mockRepository.EXPECT().CreatePost(gomock.Any()).Return(errors.New("dummy error"))

		postService := services.NewPostService(zap.NewNop(), conf, mockRepository)
		_, err := postService.CreatePost(context.Background(), &post)
		assert.NotNil(t, err)
	})
}

func TestPostService_GetPostsWithFilters(t *testing.T) {
	t.Run("given ads are enabled and post length greater than 3 when all dependecies ok "+
		"then it should return nil and should return correct response", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mocksRepository := mocks.NewMockRepositoryInterface(ctrl)
		rankedPosts := generatePosts(5)
		promotedPosts := []*models.Post{
			{
				ID:        "3",
				Title:     "Promoted Post 1",
				Author:    "t2_promoted",
				Link:      "https://example.com/promoted1",
				Subreddit: "testpromoted",
				Content:   "",
				Score:     0,
				Promoted:  false,
				NSFW:      false,
			},
		}
		conf := &config.Config{
			AdsEnabled: true,
			PostLengthAdPositionMap: map[int]int{
				3:  1,
				17: 16,
			},
			ItemPerPage: 25,
		}
		expectedResponse := &models.ListPostsResponse{
			Posts: []*models.Post{
				{
					ID:        "1",
					Title:     "Post 1",
					Author:    "t2_user123",
					Link:      "https://example.com/post1",
					Subreddit: "testsubreddit",
					Content:   "",
					Score:     1,
					Promoted:  false,
					NSFW:      false,
				},
				{
					ID:        "3",
					Title:     "Promoted Post 1",
					Author:    "t2_promoted",
					Link:      "https://example.com/promoted1",
					Subreddit: "testpromoted",
					Content:   "",
					Score:     0,
					Promoted:  false,
					NSFW:      false,
				},
				{
					ID:        "2",
					Title:     "Post 2",
					Author:    "t2_user123",
					Link:      "https://example.com/post2",
					Subreddit: "testsubreddit",
					Content:   "",
					Score:     2,
					Promoted:  false,
					NSFW:      false,
				},
				{
					ID:        "3",
					Title:     "Post 3",
					Author:    "t2_user123",
					Link:      "https://example.com/post3",
					Subreddit: "testsubreddit",
					Content:   "",
					Score:     3,
					Promoted:  false,
					NSFW:      false,
				},
				{
					ID:        "4",
					Title:     "Post 4",
					Author:    "t2_user123",
					Link:      "https://example.com/post4",
					Subreddit: "testsubreddit",
					Content:   "",
					Score:     4,
					Promoted:  false,
					NSFW:      false,
				},
				{
					ID:        "5",
					Title:     "Post 5",
					Author:    "t2_user123",
					Link:      "https://example.com/post5",
					Subreddit: "testsubreddit",
					Content:   "",
					Score:     5,
					Promoted:  false,
					NSFW:      false,
				},
			},
			Page:       1,
			TotalPages: 4,
		}
		mocksRepository.EXPECT().GetRankedPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(rankedPosts, nil)
		mocksRepository.EXPECT().GetPromotedPosts(gomock.Any()).Return(promotedPosts, nil)
		mocksRepository.EXPECT().GetTotalPostsCount(gomock.Any()).Return(100, nil)

		postService := services.NewPostService(zap.NewNop(), conf, mocksRepository)
		actualResponse, err := postService.GetPostsWithFilters(0, 25, nil)
		assert.Nil(t, err)
		assert.Equal(t, expectedResponse, actualResponse)
	})

	t.Run("given ads are disabled when all dependecies ok then it should return nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mocksRepository := mocks.NewMockRepositoryInterface(ctrl)
		rankedPosts := generatePosts(4)
		conf := &config.Config{
			AdsEnabled:  false,
			ItemPerPage: 25,
		}

		mocksRepository.EXPECT().GetRankedPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(rankedPosts, nil)
		mocksRepository.EXPECT().GetTotalPostsCount(gomock.Any()).Return(100, nil)

		postService := services.NewPostService(zap.NewNop(), conf, mocksRepository)
		_, err := postService.GetPostsWithFilters(0, 25, nil)
		assert.Nil(t, err)
	})

	t.Run("given valid request when unable to get ranked posts then it should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mocksRepository := mocks.NewMockRepositoryInterface(ctrl)
		conf := &config.Config{AdsEnabled: true}

		mocksRepository.EXPECT().GetRankedPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("dummy error"))

		postService := services.NewPostService(zap.NewNop(), conf, mocksRepository)
		_, err := postService.GetPostsWithFilters(0, 25, nil)
		assert.NotNil(t, err)
	})

	t.Run("given ads are enabled when unable to get promoted posts then it should return nil", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mocksRepository := mocks.NewMockRepositoryInterface(ctrl)
		rankedPosts := generatePosts(4)
		conf := &config.Config{
			AdsEnabled:  true,
			ItemPerPage: 25,
		}

		mocksRepository.EXPECT().GetRankedPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(rankedPosts, nil)
		mocksRepository.EXPECT().GetPromotedPosts(gomock.Any()).Return(nil, errors.New("dummy error"))
		mocksRepository.EXPECT().GetTotalPostsCount(gomock.Any()).Return(100, nil)

		postService := services.NewPostService(zap.NewNop(), conf, mocksRepository)
		_, err := postService.GetPostsWithFilters(0, 25, nil)
		assert.Nil(t, err)
	})

	t.Run("given ads are enabled when unable to get page count then it should return error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		mocksRepository := mocks.NewMockRepositoryInterface(ctrl)
		rankedPosts := generatePosts(5)
		promotedPosts := []*models.Post{
			{
				ID:        "3",
				Title:     "Promoted Post 1",
				Author:    "t2_promoted",
				Link:      "https://example.com/promoted1",
				Subreddit: "testpromoted",
				Content:   "",
				Score:     0,
				Promoted:  false,
				NSFW:      false,
			},
		}
		conf := &config.Config{AdsEnabled: true}

		mocksRepository.EXPECT().GetRankedPosts(gomock.Any(), gomock.Any(), gomock.Any()).Return(rankedPosts, nil)
		mocksRepository.EXPECT().GetPromotedPosts(gomock.Any()).Return(promotedPosts, nil)
		mocksRepository.EXPECT().GetTotalPostsCount(gomock.Any()).Return(0, errors.New("dummy error"))

		postService := services.NewPostService(zap.NewNop(), conf, mocksRepository)
		_, err := postService.GetPostsWithFilters(0, 25, nil)
		assert.NotNil(t, err)
	})
}

func generatePosts(count int) []*models.Post {
	var posts []*models.Post
	for i := 0; i < count; i++ {
		posts = append(posts, &models.Post{
			ID:        fmt.Sprint(i + 1),
			Title:     "Post " + fmt.Sprint(i+1),
			Author:    "t2_user123",
			Link:      "https://example.com/post" + fmt.Sprint(i+1),
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     i + 1,
			Promoted:  false,
			NSFW:      false,
		})
	}
	return posts
}
