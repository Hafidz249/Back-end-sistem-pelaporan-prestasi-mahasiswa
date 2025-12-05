package config

import (
	"os"
)

// GetJWTSecret mengambil JWT secret dari environment variable
func GetJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-key-change-in-production"
	}
	return secret
}

// GetDatabaseURL mengambil database URL dari environment variable
func GetDatabaseURL() string {
	return os.Getenv("DATABASE_URL")
}

// GetAppPort mengambil port aplikasi dari environment variable
func GetAppPort() string {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	return port
}