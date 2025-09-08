package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port int
}

func Load() *Config {
	port := 8080
	if portStr := os.Getenv("PORT"); portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	return &Config{
		Port: port,
	}
}