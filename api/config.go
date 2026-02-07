package app

import (
	"time"

	"github.com/krakn/expense-management-backend-go/internal/utils"
)

type ApplicationServerConfig struct {
	Version        string
	Environment    string
	Address        string
	WriteTimeout   time.Duration
	ReadTimeout    time.Duration
	MaxHeaderBytes int
}

func NewApplicationServerConfig() *ApplicationServerConfig {
	return &ApplicationServerConfig{
		Version:        "0.0.1",
		Environment:    utils.GetEnv("ENVIRONMENT", "unknown"),
		Address:        utils.GetEnv("SERVER_ADDRESS", ":8080"),
		WriteTimeout:   15 * time.Second,
		ReadTimeout:    15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

}