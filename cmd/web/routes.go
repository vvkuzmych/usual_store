package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(SessionLoad)

	mux.Get("/", app.Home)

	//mux.Post("/virtual-terminal-payment-succeeded", app.VirtualTerminalPaymentSucceeded)
	//mux.Get("/virtual-terminal-receipt", app.VirtualTerminalReceipt)

	mux.Get("/widgets/{id}", app.ChargeOnce)
	mux.Post("/payment-succeeded", app.PaymentSucceeded)
	mux.Get("/receipt", app.Receipt)

	mux.Get("/plans/golden", app.GoldenPlan)
	mux.Get("/receipt/golden", app.GoldenPlanReceipt)

	mux.Get("/login", app.LoginPage)
	mux.Post("/login", app.PostLoginPage)
	mux.Get("/logout", app.LogoutPage)

	mux.Route("/admin", func(r chi.Router) {
		r.Use(app.Auth)
		r.Get("/virtual-terminal", app.VirtualTerminal)
	})
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
