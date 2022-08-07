package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/maypok86/conduit/internal/config"
	"github.com/stretchr/testify/require"
)

type env struct {
	environment      string
	httpHost         string
	httpPort         string
	postgresHost     string
	postgresPort     string
	postgresDBName   string
	postgresUser     string
	postgresPassword string
}

func setEnv(t *testing.T, env env) {
	t.Helper()

	require.NoError(t, os.Setenv("ENVIRONMENT", env.environment))
	require.NoError(t, os.Setenv("HTTP_HOST", env.httpHost))
	require.NoError(t, os.Setenv("HTTP_PORT", env.httpPort))
	require.NoError(t, os.Setenv("POSTGRES_HOST", env.postgresHost))
	require.NoError(t, os.Setenv("POSTGRES_PORT", env.postgresPort))
	require.NoError(t, os.Setenv("POSTGRES_DBNAME", env.postgresDBName))
	require.NoError(t, os.Setenv("POSTGRES_USER", env.postgresUser))
	require.NoError(t, os.Setenv("POSTGRES_PASSWORD", env.postgresPassword))
}

func TestGet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		env  env
		want *config.Config
	}{
		{
			name: "test config",
			env: env{
				environment:      "test",
				httpHost:         "0.0.0.0",
				httpPort:         "8080",
				postgresHost:     "postgres",
				postgresPort:     "5431",
				postgresDBName:   "test_wb-l0",
				postgresUser:     "test_wb-l0",
				postgresPassword: "test",
			},
			want: &config.Config{
				Environment: "test",
				HTTP: config.HTTP{
					Host:           "0.0.0.0",
					Port:           "8080",
					MaxHeaderBytes: 1,
					ReadTimeout:    10 * time.Second,
					WriteTimeout:   10 * time.Second,
				},
				Postgres: config.Postgres{
					Host:     "postgres",
					Port:     "5431",
					DBName:   "test_wb-l0",
					User:     "test_wb-l0",
					Password: "test",
					SSLMode:  "disable",
				},
				Logger: config.Logger{
					Level: "info",
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			setEnv(t, tt.env)

			got := config.Get()
			require.True(t, reflect.DeepEqual(got, tt.want))
		})
	}
}