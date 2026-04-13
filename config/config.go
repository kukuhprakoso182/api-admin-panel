package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
	CORS     CORSConfig
}

type AppConfig struct {
	Name string
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type JWTConfig struct {
	Secret              string
	AccessExpireMinutes time.Duration
	RefreshExpireHours  time.Duration
}

type LogConfig struct {
	Dir string
}

type CORSConfig struct {
	AllowOrigins []string
	AllowHeaders []string
}

func Load() (*Config, error) {
	_ = godotenv.Load() // .env opsional — di production pakai env OS

	// ── Wajib ada — gagal keras jika tidak di-set ──────────────────────────
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET environment variable is required and cannot be empty")
	}
	if len(jwtSecret) < 32 {
		return nil, errors.New("JWT_SECRET must be at least 32 characters long")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" && os.Getenv("APP_ENV") == "production" {
		return nil, errors.New("DB_PASSWORD environment variable is required in production")
	}

	// ── Parse numerik ──────────────────────────────────────────────────────
	maxOpenConns, _       := strconv.Atoi(getEnv("DB_MAX_OPEN_CONNS", "25"))
	maxIdleConns, _       := strconv.Atoi(getEnv("DB_MAX_IDLE_CONNS", "10"))
	connMaxLifetimeSec, _ := strconv.Atoi(getEnv("DB_CONN_MAX_LIFETIME", "300"))
	accessExpireMin, _    := strconv.Atoi(getEnv("JWT_ACCESS_EXPIRES_MINUTES", "60"))
	refreshExpireHours, _ := strconv.Atoi(getEnv("JWT_REFRESH_EXPIRES_HOURS", "168"))

	// ── CORS origins ───────────────────────────────────────────────────────
	allowOrigins := strings.Split(getEnv("CORS_ALLOW_ORIGINS", "http://localhost:3000"), ",")
	allowHeaders := strings.Split(
		getEnv("CORS_ALLOW_HEADERS", "Authorization,Content-Type,X-Request-ID"),
		",",
	)

	return &Config{
		App: AppConfig{
			Name: getEnv("APP_NAME", "api-admin-panel"),
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "127.0.0.1"),
			Port:            getEnv("DB_PORT", "3306"),
			User:            getEnv("DB_USER", "root"),
			Password:        dbPassword,
			Name:            getEnv("DB_NAME", "go_service"),
			MaxOpenConns:    maxOpenConns,
			MaxIdleConns:    maxIdleConns,
			ConnMaxLifetime: time.Duration(connMaxLifetimeSec) * time.Second,
		},
		JWT: JWTConfig{
			Secret:              jwtSecret,
			AccessExpireMinutes: time.Duration(accessExpireMin) * time.Minute,
			RefreshExpireHours:  time.Duration(refreshExpireHours) * time.Hour,
		},
		Log: LogConfig{
			Dir: getEnv("LOG_DIR", "./logs"),
		},
		CORS: CORSConfig{
			AllowOrigins: allowOrigins,
			AllowHeaders: allowHeaders,
		},
	}, nil
}

func (d *DatabaseConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		d.User, d.Password, d.Host, d.Port, d.Name,
	)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
