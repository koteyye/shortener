package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
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
	Server          string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	Shortener       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:"/tmp/short-url-db.json"`
	DataBaseDns     string `env:"DATABASE_DNS"`
}

func GetConfig(logger zap.SugaredLogger) (*Config, error) {
	var flagRunAddr string
	var flagShortenerAddr string
	var flagFileStoragePath string
	var flagDataBaseDNS string

	flag.StringVar(&flagRunAddr, "a", defaultServer, "address and port to run server")
	flag.StringVar(&flagShortenerAddr, "b", defaultShortenerHost, "address and port to shortener")
	flag.StringVar(&flagFileStoragePath, "f", defaultFileStoragePath, "file path for DB")
	flag.StringVar(&flagDataBaseDNS, "d", "", "db dns")
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

	//Конфиг DB
	dbVal := calcValue(envVal.DataBaseDns, flagDataBaseDNS, "")

	cfg = &Config{
		Server: &Server{
			BaseURL: "/",
			Listen:  serverVal,
		},
		Shortener: &Shortener{
			Listen: shortenerVal,
		},
		FileStoragePath: filePathVal,
		DataBaseDNS:     dbVal,
	}

	logger.Info("Server address:", cfg.Server.Listen)
	logger.Info("BaseURL:", cfg.Shortener.Listen)
	logger.Info("File storage path:", cfg.FileStoragePath)
	logger.Info("DataBase DN:", cfg.DataBaseDNS)

	return cfg, nil
}

func calcValue(envVal string, fl string, defVal string) string {
	if envVal == defVal || fl != "" {
		return fl
	} else if envVal != defVal {
		return envVal
	} else {
		return defVal
	}
}
