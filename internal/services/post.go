package services

import (
	"math/rand"

	"broozkan/postapi/internal/models"
)

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

func prepareIndices(posts []*models.Post, adPositions map[int]int) []int {
	var adIndices []int
	for k, v := range adPositions {
		if len(posts) >= k {
			adIndices = append(adIndices, v)
		}
	}
	return adIndices
}
