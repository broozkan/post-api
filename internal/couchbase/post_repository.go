package couchbase

import (
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/couchbase/gocb/v2"
)

func NewPostRepository(cbConfig *config.Couchbase) (*Couchbase, error) {
	return New(cbConfig)
}

func (c *Couchbase) CreatePost(post *models.Post) error {
	_, err := c.PostBucket.Collection(c.PostCollection).Insert(post.ID, post, &gocb.InsertOptions{
		Expiry: 0,
	})
	if err != nil {
		return err
	}
	return nil
}
