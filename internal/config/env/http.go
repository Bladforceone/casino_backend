package env

import (
	"casino_backend/internal/config"
	"errors"
	"log"
	"net"
	"os"
)

const (
	httpHostEnvName = "HTTP_HOST"
	httpPortEnvName = "HTTP_PORT"
)

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() (config.HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		log.Printf("environment variable %s not set", httpHostEnvName)
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}
	return &httpConfig{
		host: host,
		port: port,
	}, nil
}

func (h *httpConfig) Address() string {
	return net.JoinHostPort(h.host, h.port)
}
