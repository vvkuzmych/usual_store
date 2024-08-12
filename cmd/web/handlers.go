package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func (app *application) getEnvData() map[string]string {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Get the publishable key from the environment variable
	publishableKey := os.Getenv("PUBLISHABLE_KEY")
	if publishableKey == "" {
		log.Fatalf("PUBLISHABLE_KEY not set in .env file")
	}
	stringMap := make(map[string]string)
	stringMap["publishable_key"] = publishableKey
	return stringMap
}

func (app *application) VirtualTerminal(w http.ResponseWriter, r *http.Request) {
	stringMap := app.getEnvData()

	if err := app.renderTemplate(w, r, "terminal", &templateData{
		StringMap: stringMap,
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) PaymentSucceeded(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.errorLog.Println(err)
		return
	}
	cardHolder := r.Form.Get("cardholder_name")
	email := r.Form.Get("email")
	paymentIntent := r.Form.Get("payment_intent")
	paymentMethod := r.Form.Get("payment_method")
	paymentAmount := r.Form.Get("payment_amount")
	paymentCurrency := r.Form.Get("payment_currency")

	data := make(map[string]interface{})
	data["cardholder"] = cardHolder
	data["email"] = email
	data["pi"] = paymentIntent
	data["pm"] = paymentMethod
	data["pa"] = paymentAmount
	data["pc"] = paymentCurrency

	if err := app.renderTemplate(w, r, "succeeded", &templateData{
		Data: data,
	}); err != nil {
		app.errorLog.Println(err)
	}
}

func (app *application) ChargeOnce(w http.ResponseWriter, r *http.Request) {
	if err := app.renderTemplate(w, r, "buy-once", &templateData{
		StringMap: app.getEnvData(),
	}, "stripe-js"); err != nil {
		app.errorLog.Println(err)
	}
}
