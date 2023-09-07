package config

import (
	"flag"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

type Config struct {
	Server    *Server    `yaml:"server"`
	Shortener *Shortener `yaml:"shortener"`
}

type Server struct {
	BaseURL string `yaml:"base_url"`
	Listen  string `yaml:"listen"`
}

type Shortener struct {
	BaseURL string `yaml:"base_url"`
	Listen  string `yaml:"listen"`
}

func GetConfig() (*Config, error) {
	var flagRunAddr string
	var flagShortenerAddr string
	var configPath string
	flag.StringVar(&flagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&flagShortenerAddr, "b", "", "address and port to shortener")
	flag.StringVar(&configPath, "config", "", "used for set path to config file")
	flag.Parse()
	if configPath == "" {
		configPath = "config/config.yaml"
	}

	var cfg Config
	data, err := os.ReadFile(filepath.Clean(configPath))
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	if flagRunAddr != "" {
		cfg.Server.Listen = flagRunAddr
	}
	if flagShortenerAddr != "" {
		cfg.Shortener.Listen = flagShortenerAddr
	}

	return &cfg, err
}
