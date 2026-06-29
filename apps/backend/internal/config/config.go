package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/devrapture/pod-events/pkg/crypto"
)

type Config struct {
	AppEnv                string
	Port                  string
	DatabaseURL           string
	TokenEncryptionKey    string
	SpotifyClientSecret   string
	SpotifyClientID       string
	SpotifyRedirectURL    string
	JwtExpires            int
	JwtSecret             string
	TelegramBotToken      string
	TelegramWebhookSecret string
	BotName               string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	JwtExpires, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil || JwtExpires <= 0 {
		log.Println("invalid JWT_EXPIRES_IN_HOURS, defaulting to 24")
		JwtExpires = 24
	}

	config := &Config{
		AppEnv:                getEnv("APP_ENV", "development"),
		Port:                  getEnv("PORT", "8080"),
		DatabaseURL:           mustGetEnv("DATABASE_URL"),
		TokenEncryptionKey:    mustGetEnv("TOKEN_ENCRYPTION_KEY"),
		SpotifyClientSecret:   mustGetEnv("SPOTIFY_CLIENT_SECRET"),
		SpotifyClientID:       mustGetEnv("SPOTIFY_CLIENT_ID"),
		SpotifyRedirectURL:    mustGetEnv("SPOTIFY_REDIRECT_URL"),
		JwtExpires:            JwtExpires,
		JwtSecret:             mustGetEnv("JWT_SECRET"),
		TelegramBotToken:      mustGetEnv("TELEGRAM_BOT_TOKEN"),
		TelegramWebhookSecret: mustGetEnv("TELEGRAM_WEBHOOK_SECRET"),
		BotName:               mustGetEnv("BOT_NAME"),
	}

	return config, config.validate()
}

func (c *Config) validate() error {
	if _, err := crypto.DecodeEncryptionKey(c.TokenEncryptionKey); err != nil {
		return fmt.Errorf("token encryption key: %w", err)
	}
	return nil
}

func (c *Config) IsProduction() bool {
	return c.AppEnv == "production"
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func mustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("required environment variable %q is not set", key))
	}
	return val
}
