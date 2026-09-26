package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config aggregates runtime configuration.
type Config struct {
	ServerPort    string
	CORSOrigins   string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	JWTExpire     time.Duration
	RateLimitReq  int
	RateLimitWin  time.Duration
	MinIOEndpoint string
	MinIOUser     string
	MinIOPassword string
	MinIOBucket   string
	MinIOUseSSL   bool
}

// Load reads configuration from environment variables.
func Load() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "8080"),
		CORSOrigins:   getEnv("APP_CORS_ORIGINS", "http://localhost:8012"),
		DBHost:        getEnv("DB_HOST", "db"),
		DBPort:        getEnv("DB_PORT", "5505"),
		DBUser:        getEnv("DB_USER", "gbstudyapply_user"),
		DBPassword:    getEnv("DB_PASSWORD", "gbstudyapply_pwd"),
		DBName:        getEnv("DB_NAME", "gbstudyapply_db"),
		JWTSecret:     getEnv("JWT_SECRET", "change_me_to_a_long_random_string"),
		JWTExpire:     time.Duration(getEnvInt("JWT_EXPIRE_HOURS", 72)) * time.Hour,
		RateLimitReq:  getEnvInt("RATE_LIMIT_REQUESTS", 60),
		RateLimitWin:  time.Duration(getEnvInt("RATE_LIMIT_WINDOW_SECONDS", 60)) * time.Second,
		MinIOEndpoint: getEnv("MINIO_ENDPOINT", "minio:9000"),
		MinIOUser:     getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinIOPassword: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		MinIOBucket:   getEnv("MINIO_BUCKET", "gbstudyapply"),
		MinIOUseSSL:   getEnv("MINIO_USE_SSL", "false") == "true",
	}
}

// DSN returns the PostgreSQL connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=Asia/Shanghai",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}
