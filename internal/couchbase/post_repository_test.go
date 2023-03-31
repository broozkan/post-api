package couchbase_test

import (
	"time"

	"broozkan/postapi/internal/couchbase"
	"broozkan/postapi/internal/models"

	"github.com/stretchr/testify/assert"
)

func (s *CouchbaseTestSuite) TestCreatePost() {
	s.Run("given valid post when called then it should return nil", func() {
		time.Sleep(time.Second * 2)
		repository, err := couchbase.NewPostRepository(s.couchbaseConfig)
		assert.NotNil(s.T(), repository)
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

		err = repository.CreatePost(post)
		assert.Nil(s.T(), err)
	})
}
