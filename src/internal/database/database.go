package database

import (
	"database/sql"
	"essay/src/internal/config"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	Instance *sql.DB
}

var dbInstance *DB
var once sync.Once

func GetPostgreSQLConnection() *DB {
	once.Do(func() {
		db := initializeDatabase()
		dbInstance = &DB{Instance: db}
	})
	log.Println()
	return dbInstance
}

func initializeDatabase() *sql.DB {
	cfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatal("Error loading config:", err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("Successfully connected to the database")

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db
}

func (db *DB) Close() error {
	if db.Instance != nil {
		err := db.Instance.Close()
		if err != nil {
			log.Printf("Failed to close the database connection: %v", err)
			return err
		}
		log.Println("Database connection closed.")
		db.Instance = nil
	}
	return nil
}
