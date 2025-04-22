package config

import "os"

// Config holds all configuration for the application
type Config struct {
	DatabaseURL      string
	SupabaseURL      string
	SupabaseKey      string
	StripeSecretKey  string
	StripeWebhookKey string
	JWTSecret        string
	Environment      string
	AllowedOrigins   string
}

// New creates a new Config
func New() *Config {
	return &Config{
		DatabaseURL:      getEnv("DATABASE_URL", ""),
		SupabaseURL:      getEnv("SUPABASE_URL", ""),
		SupabaseKey:      getEnv("SUPABASE_KEY", ""),
		StripeSecretKey:  getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookKey: getEnv("STRIPE_WEBHOOK_SECRET", ""),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		Environment:      getEnv("ENVIRONMENT", "development"),
		AllowedOrigins:   getEnv("ALLOWED_ORIGINS", "http://localhost:3000"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
