package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Origin", "X-Request-Id", "X-SL-TOKEN"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Route("/api/admin", func(r chi.Router) {
		mux.Use(app.Auth)

		mux.Get("/test", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("OK, got into template"))
		})

		mux.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)
	})

	mux.Post("/api/payment-intent", app.GetPaymentIntent)
	mux.Get("/api/widgets/{id}", app.GetWidgetByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is-authenticated", app.CheckAuthentication)

	return mux
}
