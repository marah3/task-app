package database

import (
	"database/sql"
	"log"
	"taskapp/config"
	"taskapp/internal/models"

	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

var DB *bun.DB

func Init(cfg *config.Config) {
	dsn := cfg.DatabaseURL
	log.Println("Using DSN:", dsn)

	sqldb, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}

	DB = bun.NewDB(sqldb, pgdialect.New())

	DB.RegisterModel((*models.TaskUser)(nil))

	if err := DB.Ping(); err != nil {
		log.Fatal("Failed to connect to DB:", err)
	}

	log.Println("Connected to PostgreSQL with Bun")
}
