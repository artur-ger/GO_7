package ledger

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func databaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}

	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	pass := getEnv("DB_PASS", "postgres")
	name := getEnv("DB_NAME", "cashapp")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, pass, host, port, name)
}

// InitDB открывает соединение с PostgreSQL через sql.Open("pgx", dsn) и проверяет Ping.
func InitDB() error {
	conn, err := sql.Open("pgx", databaseURL())
	if err != nil {
		return fmt.Errorf("не удалось открыть соединение с БД: %w", err)
	}

	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(30 * time.Minute)

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return fmt.Errorf("база данных не отвечает (проверьте DATABASE_URL / DB_*): %w", err)
	}

	db = conn
	log.Println("подключились к PostgreSQL")
	return nil
}

// Close закрывает соединение с БД и Redis.
func Close() {
	if db != nil {
		_ = db.Close()
	}
	if rdb != nil {
		_ = rdb.Close()
	}
}
