package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	ehandler "github.com/krakn/expense-management-backend-go/api/handler"
	ehandleraccount "github.com/krakn/expense-management-backend-go/api/handler/account"
	ehandlercategory "github.com/krakn/expense-management-backend-go/api/handler/category"
	ehandlertransaction "github.com/krakn/expense-management-backend-go/api/handler/transaction"
	ehandleruser "github.com/krakn/expense-management-backend-go/api/handler/user"
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
			router.Post("/login", ehandleruser.Login(a.logger, a.storage, a.authenticator))

			router.Post("/", ehandleruser.CreateUser(a.logger, a.storage))

			router.Group(func(router chi.Router) {
				router.Use(emiddleware.Authenticate(a.authenticator, a.logger, LOGGED_IN_USER_ID))

				router.Get("/{id}", ehandleruser.GetuserByID(a.logger, a.storage))
				router.Put("/", ehandleruser.UpdateUser(a.logger, a.storage, LOGGED_IN_USER_ID))
				router.Delete("/{id}", ehandleruser.DeleteUser(a.logger, a.storage, LOGGED_IN_USER_ID))
			})
		})
		router.Route("/category", func(router chi.Router) {
			router.Use(emiddleware.Authenticate(a.authenticator, a.logger, LOGGED_IN_USER_ID))

			router.Post("/", ehandlercategory.CreateCategory(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/", ehandlercategory.GetAllCategory(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/{categoryid}", ehandlercategory.GetCategoryByID(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Put("/{categoryid}", ehandlercategory.UpdateCategory(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Delete("/{categoryid}", ehandlercategory.DeleteCategory(a.logger, a.storage, LOGGED_IN_USER_ID))
		})
		router.Route("/account", func(router chi.Router) {
			router.Use(emiddleware.Authenticate(a.authenticator, a.logger, LOGGED_IN_USER_ID))

			router.Post("/", ehandleraccount.CreateAccount(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/", ehandleraccount.GetAllAccounts(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/{id}", ehandleraccount.GetAccountByID(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Put("/{id}", ehandleraccount.UpdateAccount(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Delete("/{id}", ehandleraccount.DeleteAccount(a.logger, a.storage, LOGGED_IN_USER_ID))
		})
		router.Route("/transaction", func(router chi.Router) {
			router.Use(emiddleware.Authenticate(a.authenticator, a.logger, LOGGED_IN_USER_ID))

			router.Post("/", ehandlertransaction.CreateTransaction(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/", ehandlertransaction.GetAllTransactions(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Get("/{id}", ehandlertransaction.GetTransactionByID(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Put("/{id}", ehandlertransaction.UpdateTransaction(a.logger, a.storage, LOGGED_IN_USER_ID))
			router.Delete("/{id}", ehandlertransaction.DeleteTransaction(a.logger, a.storage, LOGGED_IN_USER_ID))
		})

	})

	return router
}
