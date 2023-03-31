package handlers

import (
	"errors"
	"strings"

	"broozkan/postapi/internal/models"

	"github.com/asaskevich/govalidator"
)

func parseQueryStringParams(queryParams []byte) map[string]string {
	filterMap := make(map[string]string)
	for _, param := range queryParams {
		kv := strings.Split(string(param), "=")
		filterMap[kv[0]] = kv[1]
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
