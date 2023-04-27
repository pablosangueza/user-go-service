package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/providers/env"
)

type AppConf struct {
	LogLevel    string
	LogEncoding string
	Port        int
	// Max amount of time to wait until the service will be non-gracefully terminated
	MaxServerShutdownTime time.Duration
	// Amount of time to wait before shutdown timeout starts.
	// This is useful in kubernetes due to its eventual consistency model
	ShutdownGracePeriod time.Duration
	Tracing             struct {
		Enabled bool
		Debug   bool
	}

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     int
	// try creating database during migrations, this options is only used there
	DBCreate bool
}

const (
	envPrefix           = "DOM_USERS_SVC"
	defaultShutdownTime = 30 * time.Second
)

var k = koanf.New("_")

func Init() (*AppConf, error) {
	if err := loadDotEnv(); err != nil {
		return nil, err
	}

	err := k.Load(env.Provider(envPrefix, "_", func(s string) string {
		return strings.ToLower(strings.TrimPrefix(s, envPrefix))
	}), nil)
	if err != nil {
		return nil, err
	}

	c := defaultAppConfig()
	err = k.Unmarshal("", c)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// create new config with sane defaults to reduce the amount of local changes to get the app up and running
func defaultAppConfig() *AppConf {
	// new config with default value
	defaultCfg := &AppConf{
		LogLevel:              "debug",
		Port:                  50034,
		MaxServerShutdownTime: defaultShutdownTime,
	}
	defaultCfg.Tracing.Debug = false
	defaultCfg.DBHost = "localhost"
	defaultCfg.DBName = "golocal"
	defaultCfg.DBUser = "postgres"
	defaultCfg.DBPassword = "postgres"
	defaultCfg.DBPort = 5432

	return defaultCfg
}

// load .env file, very useful in development
func loadDotEnv() error {
	p, err := os.Executable()
	if err != nil {
		return nil
	}

	envFile := filepath.Join(filepath.Dir(p), ".env")
	if _, err := os.Stat(envFile); os.IsNotExist(err) {
		return nil
	}

	return godotenv.Load()
}
