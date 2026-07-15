package config

import "os"

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func JWTSecret() []byte {
	return []byte(getEnv("JWT_SECRET", "ecommindo-dev-secret-change-me"))
}

func Port() string {
	return getEnv("PORT", "8080")
}

func StaticDir() string {
	return getEnv("STATIC_DIR", "web")
}
