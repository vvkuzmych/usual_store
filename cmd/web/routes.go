package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)

	mux.Get("/widgets/{id}", app.ChargeOnce)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)

	mux.Get("/plans/golden", app.GoldenPlan)
	mux.Get("/receipt/golden", app.GoldenPlanReceipt)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.LogoutPage)
	mux.Get("/forgot-password", app.ForgotPassword)
	mux.Get("/reset-password", app.ShowResetPassword)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(app.Auth)
		r.Get("/virtual-terminal", app.VirtualTerminal)
		r.Get("/all-subscriptions", app.AllSubscriptions)
		r.Get("/all-sales", app.AllSales)
		r.Get("/sales/{id}", app.ShowSale)
		r.Get("/subscriptions/{id}", app.ShowSubscription)
	})
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
