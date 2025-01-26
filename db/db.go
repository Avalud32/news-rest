package db

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"news-rest/config"
	"sync"
)

var db *sqlx.DB
var once sync.Once

func GetDB(config config.Config) *sqlx.DB {
	once.Do(func() {
		var err error
		db, err = sqlx.Connect("postgres", config.GetDSN())
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
		log.Println("Connected to database")
	})
	return db
}

func PingDB(config config.Config) error {
	conn, err := sql.Open("postgres", config.GetDSN())
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}
