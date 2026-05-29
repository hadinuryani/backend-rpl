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

	// Auto-create settings (pengaturan) table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS pengaturan (
		kunci VARCHAR(100) PRIMARY KEY,
		nilai VARCHAR(255) NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create settings table: %w", err)
	}

	// Seed default settings
	settingsToSeed := map[string]string{
		"waktu_pengingat": "08:00",
		"nama_klinik":      "Klinik Indah Care Plus (IC+)",
		"alamat_klinik":    "Jl. Indah Care No. 45, Jakarta",
		"jam_kontrol":      "08:00 - selesai",
	}

	for k, v := range settingsToSeed {
		_, err = db.Exec(`INSERT IGNORE INTO pengaturan (kunci, nilai) VALUES (?, ?)`, k, v)
		if err != nil {
			return fmt.Errorf("failed to seed setting %s: %w", k, err)
		}
	}

	// Add reset_otp columns if they do not exist
	var count int
	_ = db.QueryRow(`SELECT COUNT(*) FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'users' AND COLUMN_NAME = 'reset_otp_code'`).Scan(&count)
	if count == 0 {
		_, err = db.Exec(`ALTER TABLE users ADD COLUMN reset_otp_code VARCHAR(6) NULL, ADD COLUMN reset_otp_expired_at DATETIME NULL`)
		if err != nil {
			log.Printf("Warning: failed to add reset_otp columns: %v", err)
		}
	}

	// Add is_hamil column if it does not exist
	var countHamil int
	_ = db.QueryRow(`SELECT COUNT(*) FROM information_schema.COLUMNS 
		WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'pasien' AND COLUMN_NAME = 'is_hamil'`).Scan(&countHamil)
	if countHamil == 0 {
		_, err = db.Exec(`ALTER TABLE pasien ADD COLUMN is_hamil TINYINT(1) DEFAULT 0`)
		if err != nil {
			log.Printf("Warning: failed to add is_hamil column: %v", err)
		} else {
			log.Println("Added is_hamil column to pasien table")
		}
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
