package services

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"broozkan/postapi/internal/models"
)

func AddPromotedPost(posts []*models.Post, promoted *models.Post, index int) ([]*models.Post, error) {
	if index > len(posts) {
		return nil, fmt.Errorf("invalid index %d, posts has length %d", index, len(posts))
	}

	if index == 0 {
		return append([]*models.Post{promoted}, posts...), nil
	}

	newPosts := append([]*models.Post{}, posts[:index]...)
	newPosts = append(newPosts, promoted)
	newPosts = append(newPosts, posts[index:]...)
	return newPosts, nil
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[generateRandomNumber(int64(len(letters)))]
	}
	return string(result)
}

func PrepareIndices(posts []*models.Post, adPositions map[int]int) []int {
	var adIndices []int
	for k, v := range adPositions {
		if len(posts) < k {
			continue
		}
		adIndex := v
		if adIndex == 0 || adIndex == len(posts) {
			continue
		}
		if adIndex > 0 && posts[adIndex-1].NSFW {
			continue
		}
		if adIndex < len(posts)-1 && posts[adIndex+1].NSFW {
			continue
		}
		adIndices = append(adIndices, adIndex)
	}
	return adIndices
}

func generateRandomNumber(l int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(l))
	if err != nil {
		panic(err)
	}
	return nBig.Int64()
}
