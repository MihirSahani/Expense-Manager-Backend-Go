package app

import (
	"net/http"
	"time"

	"github.com/krakn/expense-management-backend-go/api/application-server/handler"
	"github.com/krakn/expense-management-backend-go/internal/utils/authenticator"
)

type application struct {
	config *ApplicationServerConfig
	logger *logger
	authenticator authenticator.Authenticator
}

func NewApplication(config *ApplicationServerConfig) *application {
	return &application{
		config: config,
	}
}

func NewApplicationServer() *application {
	return NewApplication(NewApplicationServerConfig())
}

func (a *application) Run() {
	router := a.getRouter()
	server := &http.Server{
		Addr: a.config.Address,
		Handler: router,
		WriteTimeout: time.Duration(a.config.WriteTimeout) * time.Second,
		ReadTimeout: time.Duration(a.config.ReadTimeout) * time.Second,
		MaxHeaderBytes: a.config.MaxHeaderBytes,
	}
	a.logger.Info("Starting server")
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}

func (a *application) errorJSON(w http.ResponseWriter, status int, err error) {
	handler.WriteJSON(w, status, nil)
	a.logger.Error(err.Error())
}