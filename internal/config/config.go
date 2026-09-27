package config

import (
	"flag"
	"os"
)

// Config holds the application configuration.
type Config struct {
	Addr string
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Addr: getEnv("UNITKIT_ADDR", ":8080"),
	}
}

// Load loads configuration from defaults, env vars, and flags.
func Load() *Config {
	cfg := Default()

	flag.StringVar(&cfg.Addr, "addr", cfg.Addr, "listen address (env: UNITKIT_ADDR)")
	flag.Parse()

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
