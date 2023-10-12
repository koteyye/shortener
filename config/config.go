package config

import (
	"github.com/caarlos0/env/v6"
)

const (
	defaultServer          = "localhost:8080"
	defaultShortenerHost   = "http://localhost:8080"
	defaultFileStoragePath = "/tmp/short-url-db.json"
)

type Config struct {
	Server          *Server
	Shortener       *Shortener
	FileStoragePath string
	DataBaseDNS     string
}

type Server struct {
	BaseURL string `default:"/"`
	Listen  string
}

type Shortener struct {
	Listen string
}

type ENVValue struct {
	Server          string `env:"SERVER_ADDRESS"`
	Shortener       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDNS     string `env:"DATABASE_DNS"`
}

func GetConfig() (*Config, error) {

	var envVal ENVValue
	if err := env.Parse(&envVal); err != nil {
		return nil, err
	}

	cfg := mapEnvFlagToConfig(&envVal)

	return cfg, nil
}

func mapEnvFlagToConfig(envVal *ENVValue) *Config {
	return &Config{
		Server: &Server{
			Listen:  calcVal(envVal.Server, defaultServer),
			BaseURL: "/",
		},
		Shortener:       &Shortener{Listen: calcVal(envVal.Shortener, defaultShortenerHost)},
		FileStoragePath: calcVal(envVal.FileStoragePath, defaultFileStoragePath),
		DataBaseDNS:     calcVal(envVal.DataBaseDNS, ""),
	}

}

func calcVal(env string, def string) string {
	if env != "" {
		return env
	}
	return def
}
