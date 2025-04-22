package database

import (
	"database/sql"
	"fmt"
	"log"
	"taskapp/config"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var DB *bun.DB

func Init(cfg *config.Config) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	DB = bun.NewDB(sqldb, pgdialect.New())

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL with Bun")
}
