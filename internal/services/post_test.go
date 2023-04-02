package services_test

import (
	"testing"

	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestAddPromotedPost(t *testing.T) {
	posts := []*models.Post{
		{ID: "1", Title: "Post 1"},
		{ID: "2", Title: "Post 2"},
		{ID: "3", Title: "Post 3"},
	}

	promoted := &models.Post{ID: "4", Title: "Promoted Post"}

	t.Run("test valid index", func(t *testing.T) {
		newPosts, err := services.AddPromotedPost(posts, promoted, 1)
		assert.Nil(t, err)
		assert.Equal(t, 4, len(newPosts))
		assert.Equal(t, "1", newPosts[0].ID)
		assert.Equal(t, "Promoted Post", newPosts[1].Title)
		assert.Equal(t, "2", newPosts[2].ID)
		assert.Equal(t, "3", newPosts[3].ID)
	})

	t.Run("test index at start", func(t *testing.T) {
		newPosts, err := services.AddPromotedPost(posts, promoted, 0)
		assert.Nil(t, err)
		assert.Equal(t, 4, len(newPosts))
		assert.Equal(t, "Promoted Post", newPosts[0].Title)
		assert.Equal(t, "1", newPosts[1].ID)
		assert.Equal(t, "2", newPosts[2].ID)
		assert.Equal(t, "3", newPosts[3].ID)
	})

	t.Run("test index greater than length", func(t *testing.T) {
		_, err := services.AddPromotedPost(posts, promoted, 4)
		assert.Error(t, err)
	})
}

func TestPrepareIndices(t *testing.T) {
	posts := []*models.Post{
		{
			ID:        "1",
			Title:     "Post 1",
			Author:    "t2_user123",
			Link:      "https://example.com/post1",
			Subreddit: "testsubreddit",
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
			Score:     0,
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
			Score:     0,
			Promoted:  false,
			NSFW:      true,
		},
		{
			ID:        "4",
			Title:     "Post 4",
			Author:    "t2_user123",
			Link:      "https://example.com/post4",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     0,
			Promoted:  false,
			NSFW:      false,
		},
	}

	tests := []struct {
		name         string
		adPositions  map[int]int
		expectedList []int
	}{
		{
			name: "Valid indices",
			adPositions: map[int]int{
				2: 1,
				3: 2,
			},
			expectedList: []int{2},
		},
		{
			name: "NSFW adjacent",
			adPositions: map[int]int{
				2: 1,
				3: 2,
				4: 3,
			},
			expectedList: []int{2},
		},
		{
			name: "Index out of range",
			adPositions: map[int]int{
				10: 1,
				20: 2,
			},
			expectedList: []int(nil),
		},
		{
			name: "Start and end positions",
			adPositions: map[int]int{
				2: 0,
				3: 3,
			},
			expectedList: []int(nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			adIndices := services.PrepareIndices(posts, tt.adPositions)
			assert.Equal(t, tt.expectedList, adIndices)
		})
	}
}
