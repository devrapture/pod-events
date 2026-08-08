package config

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"

	"github.com/devrapture/pod-events/pkg/crypto"
)

type Config struct {
	AppEnv                 string
	Port                   string
	DatabaseURL            string
	TokenEncryptionKey     string
	SpotifyClientSecret    string
	SpotifyClientID        string
	SpotifyRedirectURL     string
	FrontendURL            string
	JwtExpires             int
	JwtSecret              string
	TelegramBotToken       string
	TelegramWebhookSecret  string
	BotName                string
	CronSecret             string
	SentryDSN              string
	SentryRelease          string
	SentryTracesSampleRate float64
	SentryEnableLogs       bool
	// TrustedProxies are CIDRs/IPs allowed to set X-Forwarded-For for rate limiting.
	// Empty means X-Forwarded-For is never trusted.
	TrustedProxies []*net.IPNet
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	JwtExpires, err := strconv.Atoi(getEnv("JWT_EXPIRES_IN_HOURS", "24"))
	if err != nil || JwtExpires <= 0 {
		log.Println("invalid JWT_EXPIRES_IN_HOURS, defaulting to 24")
		JwtExpires = 24
	}
	sentryTracesSampleRate, err := strconv.ParseFloat(getEnv("SENTRY_TRACES_SAMPLE_RATE", "0.1"), 64)
	if err != nil || sentryTracesSampleRate < 0 || sentryTracesSampleRate > 1 {
		return nil, fmt.Errorf("SENTRY_TRACES_SAMPLE_RATE must be a number between 0 and 1")
	}
	sentryEnableLogs, err := strconv.ParseBool(getEnv("SENTRY_ENABLE_LOGS", "true"))
	if err != nil {
		return nil, fmt.Errorf("SENTRY_ENABLE_LOGS must be true or false")
	}
	trustedProxies, err := parseTrustedProxies(getEnv("TRUSTED_PROXIES", ""))
	if err != nil {
		return nil, fmt.Errorf("TRUSTED_PROXIES: %w", err)
	}

	config := &Config{
		AppEnv:                 getEnv("APP_ENV", "development"),
		Port:                   getEnv("PORT", "8080"),
		DatabaseURL:            mustGetEnv("DATABASE_URL"),
		TokenEncryptionKey:     mustGetEnv("TOKEN_ENCRYPTION_KEY"),
		SpotifyClientSecret:    mustGetEnv("SPOTIFY_CLIENT_SECRET"),
		SpotifyClientID:        mustGetEnv("SPOTIFY_CLIENT_ID"),
		SpotifyRedirectURL:     mustGetEnv("SPOTIFY_REDIRECT_URL"),
		FrontendURL:            mustGetEnv("FRONTEND_URL"),
		JwtExpires:             JwtExpires,
		JwtSecret:              mustGetEnv("JWT_SECRET"),
		TelegramBotToken:       mustGetEnv("TELEGRAM_BOT_TOKEN"),
		TelegramWebhookSecret:  mustGetEnv("TELEGRAM_WEBHOOK_SECRET"),
		BotName:                mustGetEnv("BOT_NAME"),
		CronSecret:             mustGetEnv("CRON_SECRET"),
		SentryDSN:              getEnv("SENTRY_DSN", ""),
		SentryRelease:          getEnv("SENTRY_RELEASE", ""),
		SentryTracesSampleRate: sentryTracesSampleRate,
		SentryEnableLogs:       sentryEnableLogs,
		TrustedProxies:         trustedProxies,
	}

	return config, config.validate()
}

// parseTrustedProxies parses a comma-separated list of IPs or CIDRs into IPNet entries.
// A bare IP is treated as a single-host prefix (/32 or /128).
func parseTrustedProxies(raw string) ([]*net.IPNet, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}

	parts := strings.Split(raw, ",")
	nets := make([]*net.IPNet, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "/") {
			_, ipNet, err := net.ParseCIDR(part)
			if err != nil {
				return nil, fmt.Errorf("invalid CIDR %q: %w", part, err)
			}
			nets = append(nets, ipNet)
			continue
		}
		ip := net.ParseIP(part)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP %q", part)
		}
		bits := 128
		if ip.To4() != nil {
			bits = 32
			ip = ip.To4()
		}
		nets = append(nets, &net.IPNet{IP: ip, Mask: net.CIDRMask(bits, bits)})
	}
	return nets, nil
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
