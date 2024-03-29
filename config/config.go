package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

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
	EnbaleHTTPS     bool
	TrustSubnet     string
	GRPCServer      string
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
	EnbaleHTTPS     bool   `env:"ENABLE_HTTPS"`
	ConfigPath      string `env:"CONFIG"`
	TrustSubnet     string `env:"TRUST_SUBNET"`
	GRPCPort        string `env:"GRPC_PORT"`
}

// cliFlag флаги командной строки.
type cliFlag struct {
	flagJWT         string
	flagAddress     string
	flagShorten     string
	flagFilePath    string
	flagDSN         string
	flagPprof       string
	flagHTTPS       bool
	flagConfig      string
	flagTrustSubnet string
	flagGRPCPort    string
}

func initFlags() *cliFlag {
	cliFlags := &cliFlag{}
	if isFlagPassed("a") {
		flag.StringVar(&cliFlags.flagAddress, "a", "", "server address flag")
	}
	if isFlagPassed("b") {
		flag.StringVar(&cliFlags.flagShorten, "b", "", "shorten URL")
	}
	if isFlagPassed("f") {
		flag.StringVar(&cliFlags.flagFilePath, "f", "", "file path")
	}
	if isFlagPassed("d") {
		flag.StringVar(&cliFlags.flagDSN, "d", "", "dsn")
	}
	if isFlagPassed("j") {
		flag.StringVar(&cliFlags.flagJWT, "j", "", "jwt secret key")
	}
	if isFlagPassed("p") {
		flag.StringVar(&cliFlags.flagPprof, "p", "", "pprof address")
	}
	if isFlagPassed("s") {
		flag.BoolVar(&cliFlags.flagHTTPS, "s", false, "https")
	}
	if isFlagPassed("c") {
		flag.StringVar(&cliFlags.flagConfig, "c", "", "config path")
	}
	if isFlagPassed("config") {
		flag.StringVar(&cliFlags.flagConfig, "config", "", "config path")
	}
	if isFlagPassed("t") {
		flag.StringVar(&cliFlags.flagTrustSubnet, "t", "", "trusted subnt")
	}
	if isFlagPassed("g") {
		flag.StringVar(&cliFlags.flagGRPCPort, "g", "", "GRPC port")
	}
	flag.Parse()
	return cliFlags
}

func isFlagPassed(name string) bool {
	for _, arg := range os.Args {
		if strings.Index(arg, "-"+name) == 0 && (strings.Index(arg, "=") == len("-"+name) || strings.Index(arg, "=") == 0) {
			return true
		}
	}
	return false
}

type fileConfig struct {
	Server          string `json:"server_address"`
	Shortener       string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DataBaseDSN     string `json:"database_dsn"`
	JWTSecretKey    string `json:"JWTSecretKey"`
	Pprof           string `json:"pprof"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustSubnet     string `json:"trusted_subnet"`
	GRPCPort        string `json:"grpc_port"`
}

// ConfigFromFile получить конфиг из файла
func (c *fileConfig) ConfigFromFile(filepath string) error {
	file, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("can't open config file: %w", err)
	}
	err = json.Unmarshal(file, &c)
	if err != nil {
		return fmt.Errorf("can't unmarshal config file: %w", err)
	}
	return nil
}

// CIDR получение *IPNet из конфига
func (c *Config) CIDR() (*net.IPNet, error) {
	_, subnet, err := net.ParseCIDR(c.TrustSubnet)
	return subnet, err
}

// GetConfig получить конфигурацию приложения
func GetConfig() (*Config, error) {
	cliFlags := initFlags()

	var envVal ENVValue
	if err := env.Parse(&envVal); err != nil {
		return nil, err
	}

	var fileVal fileConfig
	if envVal.ConfigPath != "" {
		fileVal.ConfigFromFile(envVal.ConfigPath)
	}
	if cliFlags.flagConfig != "" {
		fileVal.ConfigFromFile(cliFlags.flagConfig)
	}

	cfg := mapEnvFlagToConfig(&envVal, cliFlags, &fileVal)

	return cfg, nil
}

func mapEnvFlagToConfig(envVal *ENVValue, cliFlags *cliFlag, fileVal *fileConfig) *Config {
	return &Config{
		Server: &Server{
			Listen:  calcVal(envVal.Server, cliFlags.flagAddress, fileVal.Server, defaultServer),
			BaseURL: "/",
		},
		Shortener:       &Shortener{Listen: calcVal(envVal.Shortener, cliFlags.flagShorten, fileVal.Shortener, defaultShortenerHost)},
		FileStoragePath: calcVal(envVal.FileStoragePath, cliFlags.flagFilePath, fileVal.FileStoragePath, defaultFileStoragePath),
		DataBaseDSN:     calcVal(envVal.DataBaseDSN, cliFlags.flagDSN, fileVal.DataBaseDSN, ""),
		JWTSecretKey:    calcVal(envVal.JWTSecretKey, cliFlags.flagJWT, fileVal.JWTSecretKey, deafultSecretKey),
		Pprof:           calcVal(envVal.Pprof, cliFlags.flagPprof, fileVal.Pprof, ""),
		EnbaleHTTPS:     calcHTTPS(envVal.EnbaleHTTPS, cliFlags.flagHTTPS, fileVal.EnableHTTPS),
		TrustSubnet:     calcVal(envVal.TrustSubnet, cliFlags.flagTrustSubnet, fileVal.TrustSubnet, ""),
		GRPCServer:      calcVal(envVal.GRPCPort, cliFlags.flagGRPCPort, fileVal.GRPCPort, ""),
	}

}

func calcVal(env string, fl string, configFile string, def string) string {
	if env != "" {
		return env
	}
	if fl != "" {
		return fl
	}
	if configFile != "" {
		return configFile
	}
	return def
}

func calcHTTPS(env bool, fl bool, configFile bool) bool {
	if env {
		return true
	}
	if fl {
		return true
	}
	if configFile {
		return true
	}
	return false
}
