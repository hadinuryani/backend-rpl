package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

// Config holds all configuration values loaded from environment variables.
type Config struct {
	AppHost        string
	AppPort        string
	AppEnv         string
	DatabaseURL    string
	JWTSecret      string
	JWTExpireHours int
	WAGatewayURL   string
	WAAPIToken     string
	AllowedOrigins []string
}

var (
	config *Config
	once   sync.Once
)

// GetConfig returns the singleton Config instance, loading from .env on first call.
func GetConfig() *Config {
	once.Do(func() {
		err := godotenv.Load()
		if err != nil {
			log.Println("Warning: .env file not found, using system environment variables")
		}

		expireHours, err := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
		if err != nil {
			expireHours = 24
		}

		origins := strings.Split(getEnv("ALLOWED_ORIGINS", "http://localhost:5173,http://localhost:3000"), ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
		}

		dbUser := getEnv("DB_USER", "root")
		dbPass := getEnv("DB_PASSWORD", "")
		dbHost := getEnv("DB_HOST", "localhost")
		dbPort := getEnv("DB_PORT", "3306")
		dbName := getEnv("DB_NAME", "klinik_ic")

		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true&loc=Local"
		}

		config = &Config{
			AppHost: getEnv("APP_HOST", "localhost"),
			AppPort: getEnv("APP_PORT", "8080"),
			AppEnv:  getEnv("APP_ENV", "development"),

			DatabaseURL: dbURL,

			JWTSecret:      getEnv("JWT_SECRET", "default-secret-change-me"),
			JWTExpireHours: expireHours,

			WAGatewayURL: getEnv("WA_GATEWAY_URL", "https://api.fonnte.com/send"),
			WAAPIToken:   getEnv("WA_API_TOKEN", ""),

			AllowedOrigins: origins,
		}
	})

	return config
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
