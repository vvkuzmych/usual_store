package main

import (
	"net/http"
	// "os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	// "go.opentelemetry.io/contrib/instrumentation/github.com/go-chi/chi/v5/otelchi"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// Add OpenTelemetry middleware if enabled (temporarily disabled for certificate issues)
	// if os.Getenv("OTEL_ENABLED") == "true" {
	// 	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	// 	if serviceName == "" {
	// 		serviceName = "usual-store-api"
	// 	}
	// 	mux.Use(otelchi.Middleware(serviceName, otelchi.WithChiRoutes(mux)))
	// }

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
	mux.Get("/api/widgets", app.GetAllWidgets)
	mux.Get("/api/products", app.GetAllWidgets) // Alias for /api/widgets
	mux.Get("/api/widgets/{id}", app.GetWidgetByID)
	mux.Get("/api/product/{id}", app.GetWidgetByID) // Alias for /api/widgets/{id}
	mux.Post("/api/create-customer-and-subscribe-to-plan", app.CreateCustomerAndSubscribeToPlan)

	mux.Post("/api/authenticate", app.CreateAuthToken)
	mux.Post("/api/is-authenticated", app.CheckAuthentication)
	mux.Post("/api/forgot-password", app.SendPasswordResetEmail)
	mux.Post("/api/reset-password", app.ResetPassword)
	mux.Post("/api/users", app.CreateUser)
	mux.Get("/api/users", app.GetAllUsers)
	mux.Delete("/api/users/{id}", app.DeleteUserByID)

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
