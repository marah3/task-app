package config

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Port       string `env:"PORT,default=8080"`
	DBHost     string `env:"DB_HOST,default=localhost"`
	DBPort     string `env:"DB_PORT,default=5432"`
	DBUser     string `env:"DB_USER,default=postgres"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBName     string `env:"DB_NAME,default=task_app"`
	RedisAddr  string `env:"REDIS_ADDR,default=localhost:6379"`
}

func LoadConfig() *Config {
	var cfg Config
	ctx := context.Background()

	if err := envconfig.Process(ctx, &cfg); err != nil {
		log.Fatalf("Failed to load env config: %v", err)
	}

	return &cfg
}
