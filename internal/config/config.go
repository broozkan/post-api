package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		AppName                 string
		LogLevel                string
		Server                  Server
		AdsEnabled              bool
		MinPostLengthForAd      int
		PostLengthAdPositionMap map[int]int // the key of that map represents the length of posts. If posts length greater or equal then the promoted post will insert to the corresponding value
		ItemPerPage             int
		AuthorPrefix            string
		AuthorIDLength          int
		Couchbase               Couchbase
	}

	Server struct {
		Port string
	}
)

func New(configPath, configName string) (*Config, error) {
	viper.AddConfigPath(configPath)
	viper.SetConfigName(configName)

	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	config := &Config{}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
