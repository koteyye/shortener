package config

import (
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
	})

}