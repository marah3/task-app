package config

import (
	"context"
	"github.com/joho/godotenv"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port        string `env:"PORT,default=8080"`
	DatabaseURL string `env:"DATABASE_URL"`
	RedisAddr   string `env:"REDIS_ADDR,default=localhost:6379"`
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	var cfg Config
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}

	return &cfg
}
