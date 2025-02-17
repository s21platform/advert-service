package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Platform Platform
}

type Service struct {
	Port string `env:"ADVERT_SERVICE_PORT"`
	Name string `env:"ADVERT_SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"ADVERT_SERVICE_POSTGRES_USER"`
	Password string `env:"ADVERT_SERVICE_POSTGRES_PASSWORD"`
	Database string `env:"ADVERT_SERVICE_POSTGRES_DB"`
	Host     string `env:"ADVERT_SERVICE_POSTGRES_HOST"`
	Port     string `env:"ADVERT_SERVICE_POSTGRES_PORT"`
}

type Metrics struct {
	Host string `env:"GRAFANA_HOST"`
	Port int    `env:"GRAFANA_PORT"`
}

type Platform struct {
	Env string `env:"ENV"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
