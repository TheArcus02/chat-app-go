package config

import (
	"os"
	"strconv"
)

type Config struct {
	Host string
	Port int
}

func LoadConfig() Config {
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = "127.0.0.1"
	}

	portStr := os.Getenv("SERVER_PORT")
	port := 8080
	if portStr != "" {
		if parsedPort, err := strconv.Atoi(portStr); err == nil {
			port = parsedPort
		}
	}

	return Config{
		Host: host,
		Port: port,
	}
}
