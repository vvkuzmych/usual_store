package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Add TraceMiddleware as the first middleware
	mux.Use(TraceMiddleware)
	// Add RateLimitMiddleware
	mux.Use(RateLimitMiddleware)

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	mux.Post("/api/payment-intent", app.GetPaymentIntent)
	mux.Get("/api/widgets/{id}", app.GetWidgetByID)
	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is-authenticated", app.CheckAuthentication)
	mux.Post("/api/forgot-password", app.SendPasswordResetEmail)
	mux.Post("/api/reset-password", app.ResetPassword)

	// Admin routes with authentication middleware
	mux.Route("/api/admin", func(r chi.Router) {
		r.Use(app.Auth) // Apply Auth middleware to this subrouter
		r.Post("/virtual-terminal-succeeded", app.VirtualTerminalPaymentSucceeded)
		r.Post("/all-sales", app.AllSales)
		r.Post("/all-subscriptions", app.AllSubscriptions)
		r.Post("/get-sale/{id}", app.GetSale)
		r.Post("/get-subscription/{id}", app.GetSale)
		r.Post("/refund", app.RefundCharge)
		r.Post("/cancel-subscription", app.CancelSubscription)
		r.Post("/all-users", app.AllUsers)
		r.Post("/all-users/{id}", app.ShowUser)
		r.Post("/all-users/{id}", app.ShowUser)
		r.Post("/all-users/edit/{id}", app.EditUser)
		r.Post("/all-users/delete/{id}", app.DeleteUser)
	})
	return mux
}
