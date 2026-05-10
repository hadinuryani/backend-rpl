package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Connect initializes the MySQL connection pool using the given DSN (Data Source Name).
func Connect(dsn string) error {
	if dsn == "" {
		return fmt.Errorf("DATABASE_URL (DSN) is empty")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Database connected successfully (MySQL)")
	return nil
}

// GetDB returns the current database connection pool.
func GetDB() *sql.DB {
	return db
}

// Close gracefully shuts down the connection pool.
func Close() {
	if db != nil {
		db.Close()
		log.Println("🔌 Database connection closed")
	}
}
