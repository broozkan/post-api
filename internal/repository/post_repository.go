package repository

import (
	"fmt"
	"strconv"

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
	query := "SELECT * FROM `post`.`_default`.`posts` as postData WHERE 1=1"

	if subreddit, ok := params["subreddit"]; ok {
		query += fmt.Sprintf(" AND `subreddit` = '%s'", subreddit)
	}

	if author, ok := params["author"]; ok {
		query += fmt.Sprintf(" AND `author` = '%s'", author)
	}

	if nsfw, ok := params["nsfw"]; ok {
		nsfwBool, err := strconv.ParseBool(nsfw)
		if err != nil {
			return nil, err
		}
		query += fmt.Sprintf(" AND `nsfw` = %t", nsfwBool)
	}

	query += fmt.Sprintf(" ORDER BY `score` DESC LIMIT %d OFFSET %d", limit, offset)

	result, err := c.Cluster.Query(query, &gocb.QueryOptions{Adhoc: true})
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for result.Next() {
		var postResult models.PostRow
		err = result.Row(&postResult)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &postResult.Post)
	}

	return posts, nil
}

func (c *Couchbase) GetPromotedPosts(count int) ([]*models.Post, error) {
	query := fmt.Sprintf("SELECT * FROM `post`.`_default`.`posts` as postData WHERE `promoted` = true LIMIT %d", count)

	result, err := c.Cluster.Query(query, &gocb.QueryOptions{Adhoc: true})
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for result.Next() {
		var postResult models.PostRow
		err = result.Row(&postResult)
		if err != nil {
			return nil, err
		}
		posts = append(posts, &postResult.Post)
	}

	return posts, nil
}

func (c *Couchbase) GetTotalPostsCount() (int, error) {
	query := "SELECT COUNT(*) as count FROM `post`.`_default`.`posts`"

	result, err := c.Cluster.Query(query, &gocb.QueryOptions{Adhoc: true})
	if err != nil {
		return 0, err
	}

	var countRow struct {
		Count int `json:"count"`
	}

	if result.Next() {
		err = result.Row(&countRow)
		if err != nil {
			return 0, err
		}
	}

	return countRow.Count, nil
}
