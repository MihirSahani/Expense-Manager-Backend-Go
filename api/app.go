package app

import (
	"net/http"
	"time"

	"github.com/krakn/expense-management-backend-go/api/logger"
	"github.com/krakn/expense-management-backend-go/internal/authenticator"
	"github.com/krakn/expense-management-backend-go/internal/authenticator/jwt"
	"github.com/krakn/expense-management-backend-go/storage"
)

type application struct {
	config        *ApplicationServerConfig
	logger        elogger.Logger
	authenticator authenticator.Authenticator
	storage       *storage.Storage
}

func NewApplication(config *ApplicationServerConfig) *application {
	s, err := storage.NewStorage()
	if err != nil {
		panic(err)
	}

	app := &application{
		config:        config,
		logger:        elogger.NewLogger(),
		authenticator: ejwt.NewJWTAuthenticator(),
		storage:       s,
	}

	return app
}

func NewApplicationServer() *application {
	config := NewApplicationServerConfig()
	return NewApplication(config)
}

func (a *application) Run() {
	server := &http.Server{
		Addr:           a.config.Address,
		Handler:        a.getRouter(),
		WriteTimeout:   time.Duration(a.config.WriteTimeout) * time.Second,
		ReadTimeout:    time.Duration(a.config.ReadTimeout) * time.Second,
		MaxHeaderBytes: a.config.MaxHeaderBytes,
	}
	a.logger.Info("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
