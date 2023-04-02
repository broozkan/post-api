package handlers

import (
	"errors"
	"net/url"
	"strings"

	"broozkan/postapi/internal/models"

	"github.com/asaskevich/govalidator"
)

func parseQueryStringParams(queryParams []byte) map[string]string {
	const secondary = 2
	filterMap := make(map[string]string)
	params := string(queryParams)
	kvPairs := strings.Split(params, "&")
	for _, kv := range kvPairs {
		parts := strings.Split(kv, "=")
		if len(parts) != secondary {
			continue
		}
		key := parts[0]
		value, err := url.QueryUnescape(parts[1])
		if err != nil {
			continue
		}
		filterMap[key] = value
	}
	return filterMap
}

func validatePost(post *models.Post) error {
	if post.Title == "" {
		return errors.New("title is required")
	}
	if post.Link != "" && !govalidator.IsURL(post.Link) {
		return errors.New("invalid link")
	}
	if post.Link != "" && post.Content != "" {
		return errors.New("a post cannot have both a link and content populated")
	}
	return nil
}
