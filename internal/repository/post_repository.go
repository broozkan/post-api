package repository

import (
	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"

	"github.com/couchbase/gocb/v2"
)

type PostRepository struct {
	conf   *config.Couchbase
	bucket *gocb.Bucket
}

func NewPostRepository(conf *config.Couchbase, bucket *gocb.Bucket) *PostRepository {
	return &PostRepository{
		conf:   conf,
		bucket: bucket,
	}
}

func (r *PostRepository) CreatePost(post *models.Post) error {
	_, err := r.bucket.Collection(r.conf.PostCollection).Insert(post.ID, post, &gocb.InsertOptions{
		Expiry: 0,
	})
	if err != nil {
		return err
	}
	return nil
}
