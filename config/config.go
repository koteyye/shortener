package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

const (
	defaultServer          = "localhost:8081"
	defaultShortenerHost   = "http://localhost:8081"
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

type cliFlag struct {
	flagAddress  string
	flagShorten  string
	flagFilePath string
	flagDNS      string
}

func GetConfig() (*Config, error) {
	cliFlags := &cliFlag{}
	flag.StringVar(&cliFlags.flagAddress, "a", "", "server address flag")
	flag.StringVar(&cliFlags.flagShorten, "b", "", "shorten URL")
	flag.StringVar(&cliFlags.flagFilePath, "f", "", "file path")
	flag.StringVar(&cliFlags.flagDNS, "d", "", "DNS")
	flag.Parse()
	fmt.Println(cliFlags)

	var envVal ENVValue
	if err := env.Parse(&envVal); err != nil {
		return nil, err
	}

	cfg := mapEnvFlagToConfig(&envVal, cliFlags)

	return cfg, nil
}

func mapEnvFlagToConfig(envVal *ENVValue, cliFlags *cliFlag) *Config {
	return &Config{
		Server: &Server{
			Listen:  calcVal(envVal.Server, cliFlags.flagAddress, defaultServer),
			BaseURL: "/",
		},
		Shortener:       &Shortener{Listen: calcVal(envVal.Shortener, cliFlags.flagShorten, defaultShortenerHost)},
		FileStoragePath: calcVal(envVal.FileStoragePath, cliFlags.flagFilePath, defaultFileStoragePath),
		DataBaseDNS:     calcVal(envVal.DataBaseDNS, cliFlags.flagDNS, ""),
	}

}

func calcVal(env string, fl string, def string) string {
	if env != "" {
		return env
	}
	if fl != "" {
		return fl
	}
	return def
}
