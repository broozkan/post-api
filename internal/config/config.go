package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		AppName      string
		Server       Server
		AdsEnabled   bool
		AdsFrequency int
		ItemPerPage  int
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
