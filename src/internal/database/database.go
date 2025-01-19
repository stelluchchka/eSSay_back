package database

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"essay/src/internal/config"

	_ "github.com/lib/pq"
)

var dbInstance *sql.DB
var once sync.Once

func GetPostgreSQLConnection() *sql.DB {
	once.Do(func() {
		dbInstance = initializeDatabase()
	})
	log.Println()
	return dbInstance
}

func initializeDatabase() *sql.DB {
	cfg, err := config.LoadConfig()
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

	dbInstance = db

	return dbInstance
}
func CloseDB() {
	if dbInstance != nil {
		err := dbInstance.Close()
		if err != nil {
			log.Printf("Failed to close the database connection: %v", err)
		} else {
			log.Println("Database connection closed.")
		}
	}
}
