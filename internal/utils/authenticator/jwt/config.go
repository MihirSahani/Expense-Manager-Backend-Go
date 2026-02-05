package jwt

import (
	"time"
	"github.com/krakn/expense-management-backend-go/internal/utils"
)

type JWTConfig struct {
	SecretKey []byte
	Issuer    string
	Audience  string
	DefaultExpiryHours time.Duration
}

func LoadAuthConfig() *JWTConfig {
	return &JWTConfig{
		SecretKey: []byte(utils.GetEnv("JWT_SECRET_KEY", "secretkey")),
		Issuer:    "Project-27",
		Audience:  "Project-27-Users",
		DefaultExpiryHours: 72 * time.Hour,
	}
}