package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds all application configuration loaded from environment variables.
type Config struct {
	ServerPort         string        `mapstructure:"SERVER_PORT"`
	DatabaseURL        string        `mapstructure:"DATABASE_URL"`
	RedisURL           string        `mapstructure:"REDIS_URL"`
	JWTSecret          string        `mapstructure:"JWT_SECRET"`
	JWTExpiry          time.Duration `mapstructure:"JWT_EXPIRY"`
	RefreshExpiry      time.Duration `mapstructure:"REFRESH_EXPIRY"`
	GitHubClientID     string        `mapstructure:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string        `mapstructure:"GITHUB_CLIENT_SECRET"`
	GitHubRedirectURI  string        `mapstructure:"GITHUB_REDIRECT_URI"`
	WebhookSecret      string        `mapstructure:"WEBHOOK_SECRET"`
	EncryptionKey      string        `mapstructure:"ENCRYPTION_KEY"`
	SIAKADBaseURL      string        `mapstructure:"SIAKAD_BASE_URL"`
	CORSOrigin         string        `mapstructure:"CORS_ORIGIN"`
	RateLimitIP        int           `mapstructure:"RATE_LIMIT_IP"`
	RateLimitUser      int           `mapstructure:"RATE_LIMIT_USER"`
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	viper.SetDefault("SERVER_PORT", "8080")
	viper.SetDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/swu_osr?sslmode=disable")
	viper.SetDefault("REDIS_URL", "redis://localhost:6379/0")
	viper.SetDefault("JWT_EXPIRY", "15m")
	viper.SetDefault("REFRESH_EXPIRY", "168h")
	viper.SetDefault("CORS_ORIGIN", "http://localhost:3000")
	viper.SetDefault("RATE_LIMIT_IP", 100)
	viper.SetDefault("RATE_LIMIT_USER", 300)
	viper.SetDefault("SIAKAD_BASE_URL", "https://smartone.smart-service.co.id")

	viper.AutomaticEnv()

	cfg := &Config{}

	cfg.ServerPort = viper.GetString("SERVER_PORT")
	cfg.DatabaseURL = viper.GetString("DATABASE_URL")
	cfg.RedisURL = viper.GetString("REDIS_URL")
	cfg.JWTSecret = viper.GetString("JWT_SECRET")
	cfg.JWTExpiry = viper.GetDuration("JWT_EXPIRY")
	cfg.RefreshExpiry = viper.GetDuration("REFRESH_EXPIRY")
	cfg.GitHubClientID = viper.GetString("GITHUB_CLIENT_ID")
	cfg.GitHubClientSecret = viper.GetString("GITHUB_CLIENT_SECRET")
	cfg.GitHubRedirectURI = viper.GetString("GITHUB_REDIRECT_URI")
	cfg.WebhookSecret = viper.GetString("WEBHOOK_SECRET")
	cfg.EncryptionKey = viper.GetString("ENCRYPTION_KEY")
	cfg.SIAKADBaseURL = viper.GetString("SIAKAD_BASE_URL")
	cfg.CORSOrigin = viper.GetString("CORS_ORIGIN")
	cfg.RateLimitIP = viper.GetInt("RATE_LIMIT_IP")
	cfg.RateLimitUser = viper.GetInt("RATE_LIMIT_USER")

	return cfg, nil
}
