package cards

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/v72/paymentmethod"
	"github.com/stripe/stripe-go/v72/sub"
)

type Card struct {
	Secret   string
	Key      string
	Currency string
}

type Transaction struct {
	TransactionStatusID int
	Amount              int
	Currency            string
	LastFour            string
	BankReturnCode      string
}

func (c *Card) Charge(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	return c.CreatePaymentIntent(currency, amount)
}

func (c *Card) CreatePaymentIntent(currency string, amount int) (*stripe.PaymentIntent, string, error) {
	stripe.Key = c.Secret
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(currency),
	}

	//params.AddMetadata("key", "value")
	pi, err := paymentintent.New(params)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return pi, "", nil
}

// GetPaymentMethod gets payment method by payment intent id
func (c *Card) GetPaymentMethod(s string) (*stripe.PaymentMethod, error) {
	stripe.Key = c.Secret
	pm, err := paymentmethod.Get(s, nil)
	if err != nil {
		return nil, err
	}
	return pm, nil
}

// RetrievePaymentIntent gets an existing payment intent by id
func (c *Card) RetrievePaymentIntent(s string) (*stripe.PaymentIntent, error) {
	stripe.Key = c.Secret
	pi, err := paymentintent.Get(s, nil)
	if err != nil {
		return nil, err
	}
	return pi, nil
}

// SubscribeToPlan subscribes a Stripe customer to a specified plan.
// Params:
// - customer: the Stripe customer object that will be subscribed to the plan.
// - plan: the ID of the Stripe subscription plan.
// - email: the email associated with the customer (not used in the function, but could be useful for logging or future functionality).
// - last4: the last four digits of the card being used.
// - cardType: the type of card being used (e.g., Visa, MasterCard).
// Returns:
// - string: the subscription ID if successful.
// - error: an error message if the subscription creation fails.
func (c *Card) SubscribeToPlan(customer *stripe.Customer, plan, email, last4, cardType string) (string, error) {
	stripeCustomerID := customer.ID
	items := []*stripe.SubscriptionItemsParams{
		{Plan: stripe.String(plan)},
	}
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items:    items,
	}
	params.AddMetadata("last_four", last4)
	params.AddMetadata("card_type", cardType)
	params.AddExpand("latest_invoice.payment_intent")
	subscription, err := sub.New(params)
	if err != nil {
		return "", err
	}
	return subscription.ID, nil
}

// CreateCustomer creates a new Stripe customer with a default payment method and email.
// Params:
// - pm: the payment method (typically a PaymentMethod ID) to be set as default for the customer.
// - email: the customer's email address.
// Returns:
// - *stripe.Customer: the created customer object if successful.
// - string: an error message (if there is one) based on the specific Stripe error code.
// - error: the error object if customer creation fails.
func (c *Card) CreateCustomer(pm, email string) (*stripe.Customer, string, error) {
	stripe.Key = c.Secret
	customerParams := &stripe.CustomerParams{
		Email:         stripe.String(email),
		PaymentMethod: stripe.String(pm),
		InvoiceSettings: &stripe.CustomerInvoiceSettingsParams{
			DefaultPaymentMethod: stripe.String(pm),
		},
	}
	customer, err := customer.New(customerParams)
	if err != nil {
		msg := ""
		if stripeErr, ok := err.(*stripe.Error); ok {
			msg = cardErrorMessage(stripeErr.Code)
		}
		return nil, msg, err
	}
	return customer, "", nil
}

// cardErrorMessage maps specific Stripe error codes related to card issues to user-friendly messages.
// Params:
// - code: the error code from Stripe.
// Returns:
// - string: a user-friendly error message explaining the issue with the card.
func cardErrorMessage(code stripe.ErrorCode) string {
	var msg = ""
	switch code {
	case stripe.ErrorCodeCardDeclined:
		msg = "Your card was declined"
	case stripe.ErrorCodeExpiredCard:
		msg = "Your card is expired"
	case stripe.ErrorCodeIncorrectCVC:
		msg = "Incorrect CVC code"
	case stripe.ErrorCodeIncorrectZip:
		msg = "Incorrect zip/postal code"
	case stripe.ErrorCodeAmountTooLarge:
		msg = "The amount is too large to charge to your card"
	case stripe.ErrorCodeAmountTooSmall:
		msg = "The amount is too small to charge to your card"
	case stripe.ErrorCodeBalanceInsufficient:
		msg = "Insufficient balance"
	case stripe.ErrorCodePostalCodeInvalid:
		msg = "Your postal code is invalid"
	default:
		msg = "Your card was declined"
	}
	return msg
}
