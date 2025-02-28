package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Port     string `env:"PORT"`
	Postgres Postgres
}

type Postgres struct {
	Host          string `env:"POSTGRES_HOST"`
	Port          string `env:"POSTGRES_PORT"`
	User          string `env:"POSTGRES_USER"`
	Password      string `env:"POSTGRES_PASSWORD"`
	DBName        string `env:"POSTGRES_DB"`
	MigrationPath string `env:"MIGRATIONS_PATH"`
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("failed to load env file: %s", err)
	}

	cfg := new(Config)

	if err := cleanenv.ReadEnv(cfg); err != nil {
		log.Fatalf("failed to read envs: %s", err)
	}

	return cfg
}
