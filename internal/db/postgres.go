package db

import (
	"database/sql"
	"fmt"

	yaml "userServer/internal/model/config/YAML"

	_ "github.com/lib/pq"
)

func NewPostgres(cfg *yaml.DBConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	//db.SetMaxOpenConns(25)
	//db.SetMaxIdleConns(25)
	//db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}
