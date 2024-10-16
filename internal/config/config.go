package config

import (
	"fmt"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	API_PORT          int    `env:"API_PORT,required"`
	ETH_NODE_URL      string `env:"ETH_NODE_URL,required"`
	DB_CONNECTION_URL string `env:"DB_CONNECTION_URL,required"`
	JWT_SECRET        string `env:"JWT_SECRET,required"`
	REDIS_URL         string `env:"REDIS_URL,required"`
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("unable to load .env file: %w", err)
	}

	cfg := Config{}
	err = env.Parse(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse environment variables: %w", err)
	}

	return &cfg, nil
}
