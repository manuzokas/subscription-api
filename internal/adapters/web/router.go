package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(subHandler *SubscriptionHandler, authHandler *AuthHandler, jwtSecret string) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.RegisterHandler)
		r.Post("/login", authHandler.LoginHandler)
	})

	r.Route("/subscriptions", func(r chi.Router) {
		r.Use(AuthMiddleware(jwtSecret))

		r.Post("/", subHandler.CreateSubscriptionHandler)
		r.Get("/{id}", subHandler.GetSubscriptionByIDHandler)
		r.Delete("/{id}", subHandler.CancelSubscriptionHandler)
	})

	return r
}
