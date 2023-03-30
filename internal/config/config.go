package config

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
	panic("implement me!")
}
