package config

import (
	"flag"
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
	flag.StringVar(&flagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&flagShortenerAddr, "b", "", "address and port to shortener")
	flag.Parse()

	cfg := &Config{
		Server: &Server{
			BaseURL: "/",
			Listen:  "localhost:8080",
		},
		Shortener: &Shortener{
			BaseURL: "/",
			Listen:  "localhost:8080",
		},
	}

	if flagRunAddr != "" {
		cfg.Server.Listen = flagRunAddr
	}
	if flagShortenerAddr != "" {
		cfg.Shortener.Listen = flagShortenerAddr
	}

	return cfg, nil
}
