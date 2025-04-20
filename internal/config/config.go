package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Metrics  Metrics
	Logger   Logger
	Kafka    Kafka
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

type Logger struct {
	Host string `env:"ADVERT_SERVICE_LOGGER_HOST"`
	Port string `env:"ADVERT_SERVICE_LOGGER_PORT"`
}

type Kafka struct {
	Host              string `env:"KAFKA_HOST"`
	Port              string `env:"KAFKA_PORT"`
	SetAttributeTopic string `env:"KAFKA_STAFF_SET_ATTRIBUTE"`
	Group             string `env:"KAFKA_ATTRIBUTE_GROUP"`
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
