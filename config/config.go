package config

import "os"

type Config struct {
	Port          string
	AdminPassword string
	DatabaseURL   string
	SecretKey     string
	CheckInterval int
}

func Load() *Config {
	c := &Config{
		Port:          getEnv("PORT", "3000"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin"),
		DatabaseURL:   getEnv("DATABASE_URL", "sqlite:///data/omega.db"),
		SecretKey:     getEnv("SECRET_KEY", "change-me-to-random"),
		CheckInterval: 60,
	}
	if v := os.Getenv("CHECK_INTERVAL"); v != "" {
		var n int
		for _, ch := range v {
			if ch >= '0' && ch <= '9' {
				n = n*10 + int(ch-'0')
			}
		}
		if n > 0 {
			c.CheckInterval = n
		}
	}
	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
