package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testServer      = "localhost:8082"
	testShortener   = "http://localhost:8082"
	testDBDSN       = "postgresql://postgres:postgres@localhost:5432/shortener?sslmode=disable"
	testConfigPath  = "./testConfig.json"
	testFileStorage = "/path/to/file.db"
	testSecretKey   = "jpoifjewf4093fgu902fj9023jf092jfc023f"
)

func TestConfig_GetConfig(t *testing.T) {
	t.Run("get config", func(t *testing.T) {
		t.Run("config file", func(t *testing.T) {
			t.Setenv("SERVER_ADDRESS", testServer)
			t.Setenv("BASE_URL", testShortener)
			t.Setenv("CONFIG", testConfigPath)

			cfg, err := GetConfig()
			wantCfg := &Config{
				Server: &Server{
					Listen:  testServer,
					BaseURL: "/",
				},
				Shortener: &Shortener{
					Listen: testShortener,
				},
				DataBaseDSN:     testDBDSN,
				FileStoragePath: testFileStorage,
				JWTSecretKey:    testSecretKey,
				EnbaleHTTPS:     true,
			}
			assert.NoError(t, err)
			assert.Equal(t, wantCfg, cfg)
		})
		t.Run("nothing", func(t *testing.T) {
			cfg, err := GetConfig()
			wantCfg := &Config{
				Server: &Server{
					Listen: defaultServer,
					BaseURL: "/",
				},
				Shortener: &Shortener{
					Listen: defaultShortenerHost,
				},
				FileStoragePath: defaultFileStoragePath,
				JWTSecretKey: deafultSecretKey,
			}
			assert.NoError(t, err)
			assert.Equal(t, wantCfg, cfg)
		})
		t.Run("flags", func(t *testing.T) {
			oldArg := os.Args
			defer func() {
				os.Args = oldArg
			}()

			args := []string{"-a", "localhost:8083"}
			os.Args = args

			cfg, err := GetConfig()
			wantCfg := &Config{
				Server: &Server{
					Listen: "localhost:8083",
					BaseURL: "/",
				},
				Shortener: &Shortener{
					Listen: defaultShortenerHost,
				},
				FileStoragePath: defaultFileStoragePath,
				JWTSecretKey: deafultSecretKey,
			}
			assert.NoError(t, err)
			assert.Equal(t, wantCfg, cfg)
		})
	})

}
