package config

import "os"

type Config struct {
DatabaseURL string
JWTSecret   string
}

func Load() Config {
return Config{
DatabaseURL: env("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=homelab port=5432 sslmode=disable"),
JWTSecret:   env("JWT_SECRET", "supersecret-change-me-in-prod"),
}
}

func env(k, def string) string {
if v := os.Getenv(k); v != "" {
return v
}
return def
}