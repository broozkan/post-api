package repository_test

import (
	"context"
	"testing"

	"broozkan/postapi/internal/config"
	"broozkan/postapi/internal/models"
	"broozkan/postapi/internal/repository"

	"github.com/couchbase/gocb/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

func TestPostRepository_CreatePost(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	t.Run("", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		req := testcontainers.ContainerRequest{
			Image:        "couchbase/server:7.0.2",
			ExposedPorts: []string{"8091/tcp", "8092/tcp", "8093/tcp", "8094/tcp", "11210/tcp"},
			Name:         "couchbase_test",
		}
		couchbaseCtnr, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			t.Fatal(err)
		}
		defer couchbaseCtnr.Terminate(ctx)

		conf := &config.Couchbase{
			ConnectionString: "couchbase://localhost",
			Username:         "admin",
			Password:         "password",
			BucketName:       "testbucket",
			PostCollection:   "posts",
		}

		opts := gocb.ClusterOptions{
			Username: conf.Username,
			Password: conf.Password,
		}
		cluster, err := gocb.Connect(conf.ConnectionString, opts)
		if err != nil {
			t.Fatal(err)
		}
		defer cluster.Close(nil)

		bucket := cluster.Bucket(conf.BucketName)

		repo := repository.NewPostRepository(conf, bucket)

		post := &models.Post{
			ID:        uuid.New().String(),
			Title:     "Test post",
			Author:    "t2_user123",
			Link:      "https://example.com",
			Subreddit: "testsubreddit",
			Content:   "Test post content",
			Score:     0,
			Promoted:  false,
			NSFW:      false,
		}

		err = repo.CreatePost(post)
		assert.Nil(t, err)
	})
}
