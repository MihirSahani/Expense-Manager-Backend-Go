package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	emiddleware "github.com/krakn/expense-management-backend-go/api/middleware"
)

func (a *application) getRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)


	router.Route("/api/"+a.config.Version, func(router chi.Router) {
		router.Get("/health", ehandler.Health(a.config.Version, a.config.Environment))

		router.Route("/user", func(router chi.Router) {
			router.Post("/login", ehandler.Login(a.logger, a.storage, a.authenticator))

			router.Post("/", ehandler.CreateUser(a.logger, a.storage))

			router.Group(func(router chi.Router) {
				router.Use(emiddleware.Authenticate(a.authenticator, a.logger, LOGGED_IN_USER_ID))

				router.Get("/{id}", ehandler.GetuserByID(a.logger, a.storage))
				router.Put("/", ehandler.UpdateUser(a.logger, a.storage, LOGGED_IN_USER_ID))
				router.Delete("/{id}", ehandler.DeleteUser(a.logger, a.storage, LOGGED_IN_USER_ID))
			})
		})

	})

	return router
}
