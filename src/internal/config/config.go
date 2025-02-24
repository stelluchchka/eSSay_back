package config

import (
	"essay/src/internal/middleware"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

type KafkaConfig struct {
	Brokers  string
	Topic    string
	ClientID string
	Acks     string
}

func LoadDBConfig() (*DBConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	config := &DBConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "1234"),
		DBName:     getEnv("DB_NAME", "db"),
	}

	return config, nil
}

func LoadKafkaConfig() (*KafkaConfig, error) {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found or failed to load")
	}

	config := &KafkaConfig{
		Brokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		Topic:    getEnv("KAFKA_TOPIC", "essay_check_queue"),
		ClientID: getEnv("KAFKA_CLIENT_ID", "essay_producer"),
		Acks:     getEnv("KAFKA_ACKS", "all"),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var SessionStore *sessions.CookieStore

func InitSessionStore() {
	SessionStore = sessions.NewCookieStore([]byte(getEnv("SECRET_KEY", "SECRET_KEYSECRET_KEY")))
	SessionStore.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   7 * 24 * 60 * 60, // неделя
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
}

var СorsConfig = middleware.CORSConfig{
	AllowedOrigins:   []string{"http://localhost:3000"},
	AllowedMethods:   []string{http.MethodDelete, http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
	AllowedHeaders:   []string{"Content-Type", "Authorization"},
	ExposeHeaders:    []string{"Content-Length"},
	AllowCredentials: true,
	MaxAge:           3600,
}
