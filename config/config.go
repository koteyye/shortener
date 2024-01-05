package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
)

const (
	defaultServer          = "localhost:8080"
	defaultShortenerHost   = "http://localhost:8080"
	defaultFileStoragePath = "/tmp/short-url-db.json"
	deafultSecretKey       = "jpoifjewf4093fgu902fj9023jf092jfc023f"
)

// Config конфигурация сервиса.
type Config struct {
	Server          *Server
	Shortener       *Shortener
	FileStoragePath string
	DataBaseDSN     string
	JWTSecretKey    string
	Pprof           string
}

// Server сервер конфигурации сервиса.
type Server struct {
	BaseURL string `default:"/"`
	Listen  string
}

// Shortener адрес сокращателя ссылок.
type Shortener struct {
	Listen string
}

// ENVValue конфигурация переменного окружения.
type ENVValue struct {
	Server          string `env:"SERVER_ADDRESS"`
	Shortener       string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DataBaseDSN     string `env:"DATABASE_DSN"`
	JWTSecretKey    string `env:"JWTSecretKey"`
	Pprof           string `env:"PPROF"`
}

// cliFlag флаги командной строки.
type cliFlag struct {
	flagJWT      string
	flagAddress  string
	flagShorten  string
	flagFilePath string
	flagDSN      string
	flagPprof    string
}

// GetConfig получение конфигурации.
func GetConfig() (*Config, error) {
	cliFlags := &cliFlag{}
	flag.StringVar(&cliFlags.flagAddress, "a", "", "server address flag")
	flag.StringVar(&cliFlags.flagShorten, "b", "", "shorten URL")
	flag.StringVar(&cliFlags.flagFilePath, "f", "", "file path")
	flag.StringVar(&cliFlags.flagDSN, "d", "", "dsn")
	flag.StringVar(&cliFlags.flagJWT, "j", "", "jwt secret key")
	flag.StringVar(&cliFlags.flagPprof, "p", "", "pprof address")
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
		DataBaseDSN:     calcVal(envVal.DataBaseDSN, cliFlags.flagDSN, ""),
		JWTSecretKey:    calcVal(envVal.JWTSecretKey, cliFlags.flagJWT, deafultSecretKey),
		Pprof:           calcVal(envVal.Pprof, cliFlags.flagPprof, ""),
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
