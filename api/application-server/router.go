package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/krakn/expense-management-backend-go/api/application-server/handler"
)

func (a *application) getRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Route("/api/v1", func(router chi.Router) {
		router.Get("/health", handler.Health)

		router.Group(func(router chi.Router) {
			router.Use(a.Authentication)
		})
	})

	return router
}