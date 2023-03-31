package repository

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

func (c *Couchbase) GetRankedPosts(offset, limit int, params map[string]string) ([]*models.Post, error) {
	panic("implement me!")
}

func (c *Couchbase) GetPromotedPosts(count int) ([]*models.Post, error) {
	panic("implement me!")
}

func (c *Couchbase) GetTotalPostsCount() (int, error) {
	panic("implement me!")
}
