package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"strings"
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
}

type Server struct {
	BaseURL string `default:"/"`
	Listen  string
}

type Shortener struct {
	Listen string
}

type ENVValue struct {
	Server          string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	Shortener       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
}

func GetConfig() (*Config, error) {
	var flagRunAddr string
	var flagShortenerAddr string
	var flagFileStoragePath string

	flag.StringVar(&flagRunAddr, "a", defaultServer, "address and port to run server")
	flag.StringVar(&flagShortenerAddr, "b", defaultShortenerHost, "address and port to shortener")
	flag.StringVar(&flagFileStoragePath, "f", defaultFileStoragePath, "file path for DB")
	flag.Parse()

	var envVal ENVValue
	if err := env.Parse(&envVal); err != nil {
		return nil, err
	}

	var cfg *Config

	//Конфиг сервера
	serverVal := calcValue(envVal.Server, flagRunAddr, defaultServer)

	//Конфиг сокращателя URL
	shortenerVal := calcValue(envVal.Shortener, flagShortenerAddr, defaultShortenerHost)

	//Конфиг файл для хранения сокращенных URL
	filePathVal := calcValue(envVal.FileStoragePath, flagFileStoragePath, defaultFileStoragePath)

	cfg = &Config{
		Server: &Server{
			BaseURL: "/",
			Listen:  serverVal,
		},
		Shortener: &Shortener{
			Listen: shortenerVal,
		},
		FileStoragePath: filePathVal,
	}

	fmt.Printf("\nServer address %v\n", cfg.Server.Listen)
	fmt.Printf("Base url %v\n", cfg.Shortener.Listen)
	fmt.Printf("File storage path %v\n", cfg.FileStoragePath)

	return cfg, nil
}

func calcValue(envVal string, fl string, defVal string) string {

	val := defVal

	if !strings.Contains(envVal, defVal) {
		val = envVal
	}
	if !strings.Contains(fl, defVal) {
		val = fl
	}

	return val
}
