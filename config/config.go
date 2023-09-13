package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Server    *Server
	Shortener *Shortener
}

type Server struct {
	BaseURL string `default:"/"`
	Listen  string
}

type Shortener struct {
	Listen string
}

type ENVValue struct {
	Server    string `env:"SERVER_ADDRESS"`
	Shortener string `env:"BASE_URL"`
}

func GetConfig() (*Config, error) {
	var flagRunAddr string
	var flagShortenerAddr string

	flag.StringVar(&flagRunAddr, "a", "", "address and port to run server")
	flag.StringVar(&flagShortenerAddr, "b", "", "address and port to shortener")
	flag.Parse()

	var envVal ENVValue
	if err := env.Parse(&envVal); err != nil {
		return nil, err
	}

	var cfg *Config

	//Конфиг сервера
	var serverVal string
	if envVal.Server != "" {
		serverVal = envVal.Server
	} else if flagRunAddr != "" {
		serverVal = flagRunAddr
	} else {
		serverVal = "localhost:8080"
	}

	//Конфиг сокращателя URL
	var shortenerVal string
	if envVal.Shortener != "" {
		shortenerVal = envVal.Shortener
	} else if flagRunAddr != "" {
		shortenerVal = flagShortenerAddr
	} else {
		shortenerVal = "http://localhost:8080"
	}

	cfg = &Config{
		Server: &Server{
			BaseURL: "/",
			Listen:  serverVal,
		},
		Shortener: &Shortener{
			Listen: shortenerVal,
		},
	}

	fmt.Printf("\nServer address %v\n", cfg.Server.Listen)
	fmt.Printf("Base url %v\n", cfg.Shortener.Listen)

	return cfg, nil
}
