package repository_test

import (
	"fmt"
	"strings"
	"time"

	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/repository"

	"github.com/couchbase/gocb/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	bucketPost      = "post"
	collectionPosts = "posts"
)

func (s *CouchbaseTestSuite) TestCreatePost() {
	s.Run("given valid post when called then it should return nil", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		post := &models.Post{
			ID:        "e64fc2df-4196-491b-9369-ee234855145d",
			Title:     "Test post",
			Author:    "t2_user123",
			Link:      "https://example.com",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     100,
			Promoted:  false,
			NSFW:      false,
		}

		err = repo.CreatePost(post)
		assert.Nil(s.T(), err)
	})
}

func (s *CouchbaseTestSuite) TestGetRankedPosts() {
	s.Run("given no documents when called then it should return error", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		s.clearCollection(repo.Cluster)
		time.Sleep(time.Second * 1)

		posts, err := repo.GetRankedPosts(0, 10, map[string]string{})

		assert.Nil(s.T(), posts)
		assert.Nil(s.T(), err)
		assert.Empty(s.T(), posts)
	})

	s.Run("given valid params when called then it should return a list of posts without error", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), repo)

		s.clearCollection(repo.Cluster)
		err = s.createPrimaryIndex(repo.Cluster, collectionPosts)
		assert.Nil(s.T(), err)
		time.Sleep(time.Second * 1)

		posts := []*models.Post{
			{
				ID:        "1",
				Title:     "Post 1",
				Author:    "author1",
				Link:      "https://example.com/post1",
				Subreddit: "testsubreddit",
				Content:   "This is the content of Post 1",
				Score:     10,
				Promoted:  false,
				NSFW:      false,
			},
			{
				ID:        "2",
				Title:     "Post 2",
				Author:    "author2",
				Link:      "https://example.com/post2",
				Subreddit: "testsubreddit",
				Content:   "This is the content of Post 2",
				Score:     20,
				Promoted:  false,
				NSFW:      false,
			},
			{
				ID:        "3",
				Title:     "Post 3",
				Author:    "author3",
				Link:      "https://example.com/post3",
				Subreddit: "testsubreddit",
				Content:   "This is the content of Post 3",
				Score:     30,
				Promoted:  false,
				NSFW:      false,
			},
		}
		for _, post := range posts {
			err = repo.CreatePost(post)
			assert.Nil(s.T(), err)
		}

		time.Sleep(time.Second * 1)

		posts, err = repo.GetRankedPosts(0, 10, map[string]string{})
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), posts)
		assert.Equal(s.T(), 3, len(posts))
		assert.Equal(s.T(), "3", posts[0].ID) // The highest score should come first
		assert.Equal(s.T(), "1", posts[2].ID) // The lowest score should come last
	})

	s.Run("given some filter params when called then it should return filtered list of posts without error", func() {
		posts := []*models.Post{
			{
				ID:        "1",
				Title:     "Post 1",
				Author:    "user1",
				Link:      "https://example.com/1",
				Subreddit: "testsub",
				Content:   "Post 1 content",
				Score:     10,
				Promoted:  false,
				NSFW:      false,
			},
			{
				ID:        "2",
				Title:     "Post 2",
				Author:    "user1",
				Link:      "https://example.com/2",
				Subreddit: "testsub",
				Content:   "Post 2 content",
				Score:     20,
				Promoted:  true,
				NSFW:      false,
			},
			{
				ID:        "3",
				Title:     "Post 3",
				Author:    "user2",
				Link:      "https://example.com/3",
				Subreddit: "testsub",
				Content:   "Post 3 content",
				Score:     30,
				Promoted:  false,
				NSFW:      true,
			},
		}
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), repo)

		s.clearCollection(repo.Cluster)
		err = s.createPrimaryIndex(repo.Cluster, collectionPosts)
		assert.Nil(s.T(), err)

		time.Sleep(time.Second * 1)

		for _, p := range posts {
			err = repo.CreatePost(p)
			require.Nil(s.T(), err)
		}
		time.Sleep(time.Second * 1)

		params := map[string]string{
			"author": "user1",
			"nsfw":   "false",
		}
		posts, err = repo.GetRankedPosts(0, 10, params)
		require.Nil(s.T(), err)
		require.Len(s.T(), posts, 1)

		require.Equal(s.T(), posts[0].ID, "1")
	})
}

func (s *CouchbaseTestSuite) TestGetPromotedPosts() {
	s.Run("given valid params when called then it should return a list of promoted posts without error", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		s.clearCollection(repo.Cluster)
		err = s.createPrimaryIndex(repo.Cluster, collectionPosts)
		assert.Nil(s.T(), err)
		time.Sleep(time.Second * 1)

		post1 := &models.Post{
			ID:        "e64fc2df-4196-491b-9369-ee234855145d",
			Title:     "Test post 1",
			Author:    "t2_user123",
			Link:      "https://example.com",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     100,
			Promoted:  true,
			NSFW:      false,
		}
		err = repo.CreatePost(post1)
		assert.Nil(s.T(), err)
		time.Sleep(time.Second * 1)

		post2 := &models.Post{
			ID:        "e64fc2df-4196-491b-9369-ee234855145e",
			Title:     "Test post 2",
			Author:    "t2_user123",
			Link:      "https://example.com",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     90,
			Promoted:  true,
			NSFW:      false,
		}
		err = repo.CreatePost(post2)
		assert.Nil(s.T(), err)
		time.Sleep(time.Second * 1)

		promotedPosts, err := repo.GetPromotedPosts(2)
		assert.Nil(s.T(), err)
		assert.NotNil(s.T(), promotedPosts)
		assert.Equal(s.T(), len(promotedPosts), 2)

		// Verify the posts are in the expected order
		assert.Equal(s.T(), promotedPosts[0].ID, post1.ID)
		assert.Equal(s.T(), promotedPosts[1].ID, post2.ID)
	})

	s.Run("given no documents when called then it should return error", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		s.clearCollection(repo.Cluster)

		posts, _ := repo.GetPromotedPosts(10)
		assert.Equal(s.T(), 0, len(posts))
	})
}

func (s *CouchbaseTestSuite) TestGetTotalCount() {
	s.Run("given zero document when called then it should return valid count", func() {
		time.Sleep(time.Second * 1)
		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		s.clearCollection(repo.Cluster)
		time.Sleep(time.Second * 1)

		count, _ := repo.GetTotalPostsCount()
		assert.Equal(s.T(), 0, count)
	})
	s.Run("given two document when called then it should return valid count", func() {
		post1 := &models.Post{
			ID:        "a1b2c3",
			Title:     "Test post 1",
			Author:    "t2_user123",
			Link:      "https://example.com/1",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     10,
			Promoted:  false,
			NSFW:      false,
		}
		post2 := &models.Post{
			ID:        "d4e5f6",
			Title:     "Test post 2",
			Author:    "t2_user456",
			Link:      "https://example.com/2",
			Subreddit: "testsubreddit",
			Content:   "",
			Score:     20,
			Promoted:  false,
			NSFW:      false,
		}

		repo, err := repository.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repo)
		assert.Nil(s.T(), err)

		s.clearCollection(repo.Cluster)
		time.Sleep(time.Second * 1)

		err = repo.CreatePost(post1)
		assert.Nil(s.T(), err)

		err = repo.CreatePost(post2)
		assert.Nil(s.T(), err)

		time.Sleep(time.Second * 1)

		count, _ := repo.GetTotalPostsCount()
		assert.Equal(s.T(), 2, count)
	})
}

func (s *CouchbaseTestSuite) clearCollection(c *gocb.Cluster) {
	_, err := c.Query(fmt.Sprintf("DELETE FROM `%s`.`_default`.`%s`", bucketPost, collectionPosts), nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	assert.Nil(s.T(), err)
}

func (s *CouchbaseTestSuite) createPrimaryIndex(c *gocb.Cluster, collection string) error {
	queryString := fmt.Sprintf("CREATE PRIMARY INDEX ON `default`:`post`.`_default`.`%s`", collection)
	_, err := c.Query(queryString, nil)
	if err != nil && strings.Contains(err.Error(), "already exists") {
		return nil
	}
	return err
}
