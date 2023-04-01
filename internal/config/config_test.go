package config_test

import (
	"testing"

	"broozkan/postapi/internal/config"

	"github.com/stretchr/testify/assert"
)

const (
	configPath = "../../test/testdata"
	configName = "test-config"
)

func TestConfig_New(t *testing.T) {
	t.Run("test given test config file when I call new then it should return config", func(t *testing.T) {
		actualConfig, _ := config.New(configPath, configName)

		expectedConfig := &config.Config{
			AppName:            "something-special",
			Server:             config.Server{Port: "1111"},
			AdsEnabled:         false,
			MinPostLengthForAd: 3,
			AdsPositions: map[int]int{
				3:  2,
				17: 16,
			},
			ItemPerPage:    27,
			AuthorPrefix:   "t2",
			AuthorIDLength: 8,
			Couchbase: config.Couchbase{
				URL:      "url",
				Username: "username",
				Password: "password",
				Buckets: []config.BucketConfig{
					{
						Name:               "post",
						CreatePrimaryIndex: false,
						Scopes: []config.ScopeConfig{{
							Name: "",
							Collections: []config.CollectionConfig{
								{
									Name:               "posts",
									CreatePrimaryIndex: true,
								},
							},
						}},
					},
				},
			},
		}

		assert.Equal(t, expectedConfig, actualConfig)
	})

	t.Run("test given non existing file when I call new then it should return error", func(t *testing.T) {
		fakeConfigPath := "../test/fake"
		notExistingConfigName := "nothing"

		_, err := config.New(fakeConfigPath, notExistingConfigName)

		assert.NotNil(t, err)
	})

	t.Run("test given bad configuration file when I call new then it should return error", func(t *testing.T) {
		badConfigName := "test-bad-config"

		_, err := config.New(configPath, badConfigName)

		assert.NotNil(t, err)
	})
}
