package config

import (
	"strings"

	"coupon_service/pkg"

	"github.com/caarlos0/env/v7"
	"github.com/joho/godotenv"
)

const (
	TestEnv        = "test"
	DevelopmentEnv = "development"
	StagingEnv     = "staging"
	ProductionEnv  = "production"
)

type Config struct {
	Env Env
}
type Env struct {
	Environment string `env:"API_ENV,notEmpty"`
	LogLevel    string `env:"API_LOG_LEVEL" envDefault:"info"`
	Port        int    `env:"API_PORT"`
	AuthConfig  struct {
		JWTSecret string `env:"JWT_SECRET,required"`
	}
}

func New() (Config, error) {
	cfgEnv, err := loadEnv()
	if err != nil {
		return Config{}, err
	}

	return Config{
		Env: cfgEnv,
	}, nil
}

func loadEnv() (Env, error) {
	err := godotenv.Load()
	if err != nil {
		return Env{}, pkg.Errorf(pkg.EINTERNAL, "failed to load .env file", err)
	}
	e := Env{}
	if err := env.Parse(&e); err != nil {
		return e, pkg.Errorf(pkg.EINTERNAL, "failed to parse env variables", err)
	}
	e.Environment = strings.ToLower(e.Environment)
	e.LogLevel = strings.ToLower(e.LogLevel)

	return e, nil
}
