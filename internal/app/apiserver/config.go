package apiserver

import (
	"fmt"
	"os"
)

// Config ...
type Config struct {
	LogLevel    string
	Port        string
	Endpoint    string
	Schema      string
	DatabaseURL string
}

// NewConfig ...
func NewConfig() *Config {

	DBHost := os.Getenv("DB_HOST")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASS")
	DB := os.Getenv("DB_NAME")
	MaxConnection := os.Getenv("PG_MAXCONN")

	dbURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable&pool_max_conns=%s", DBUser, DBPassword, DBHost, DB, MaxConnection)
	// "postgres://postgres:password@localhost/crud-db?sslmode=disable&pool_max_conns=10"

	return &Config{
		LogLevel:    "debug",
		Port:        os.Getenv("PORT"),
		Endpoint:    os.Getenv("ENDPOINT"),
		Schema:      os.Getenv("DB_SCHEMA"),
		DatabaseURL: dbURL,
	}
}
