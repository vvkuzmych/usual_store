package cards

import (
	"testing"

	"github.com/stripe/stripe-go/v72"
)

func TestCardErrorMessage(t *testing.T) {
	tests := []struct {
		name     string
		code     stripe.ErrorCode
		expected string
	}{
		{
			name:     "card declined",
			code:     stripe.ErrorCodeCardDeclined,
			expected: "Your card was declined",
		},
		{
			name:     "expired card",
			code:     stripe.ErrorCodeExpiredCard,
			expected: "Your card is expired",
		},
		{
			name:     "incorrect CVC",
			code:     stripe.ErrorCodeIncorrectCVC,
			expected: "Incorrect CVC code",
		},
		{
			name:     "incorrect ZIP",
			code:     stripe.ErrorCodeIncorrectZip,
			expected: "Incorrect zip/postal code",
		},
		{
			name:     "amount too large",
			code:     stripe.ErrorCodeAmountTooLarge,
			expected: "The amount is too large to charge to your card",
		},
		{
			name:     "amount too small",
			code:     stripe.ErrorCodeAmountTooSmall,
			expected: "The amount is too small to charge to your card",
		},
		{
			name:     "balance insufficient",
			code:     stripe.ErrorCodeBalanceInsufficient,
			expected: "Insufficient balance",
		},
		{
			name:     "postal code invalid",
			code:     stripe.ErrorCodePostalCodeInvalid,
			expected: "Your postal code is invalid",
		},
		{
			name:     "unknown error code - default message",
			code:     stripe.ErrorCode("unknown_error"),
			expected: "Your card was declined",
		},
		{
			name:     "empty error code - default message",
			code:     stripe.ErrorCode(""),
			expected: "Your card was declined",
		},
		{
			name:     "processing error - default message",
			code:     stripe.ErrorCodeProcessingError,
			expected: "Your card was declined",
		},
		{
			name:     "invalid number - default message",
			code:     stripe.ErrorCodeInvalidNumber,
			expected: "Your card was declined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cardErrorMessage(tt.code)
			if result != tt.expected {
				t.Errorf("cardErrorMessage(%v) = %q, want %q", tt.code, result, tt.expected)
			}
		})
	}
}

// Test that all expected error codes are handled
func TestCardErrorMessageCoverage(t *testing.T) {
	// List of all card-related Stripe error codes we should handle
	cardErrorCodes := []stripe.ErrorCode{
		stripe.ErrorCodeCardDeclined,
		stripe.ErrorCodeExpiredCard,
		stripe.ErrorCodeIncorrectCVC,
		stripe.ErrorCodeIncorrectZip,
		stripe.ErrorCodeAmountTooLarge,
		stripe.ErrorCodeAmountTooSmall,
		stripe.ErrorCodeBalanceInsufficient,
		stripe.ErrorCodePostalCodeInvalid,
	}

	for _, code := range cardErrorCodes {
		t.Run(string(code), func(t *testing.T) {
			msg := cardErrorMessage(code)

			// Message should not be empty
			if msg == "" {
				t.Errorf("cardErrorMessage(%v) returned empty string", code)
			}

			// Message should be user-friendly (not just the error code)
			if msg == string(code) {
				t.Errorf("cardErrorMessage(%v) returned raw error code instead of user-friendly message", code)
			}
		})
	}
}

// Benchmark for error message generation
func BenchmarkCardErrorMessage(b *testing.B) {
	codes := []stripe.ErrorCode{
		stripe.ErrorCodeCardDeclined,
		stripe.ErrorCodeExpiredCard,
		stripe.ErrorCodeIncorrectCVC,
		stripe.ErrorCodeIncorrectZip,
		"unknown_error",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code := codes[i%len(codes)]
		_ = cardErrorMessage(code)
	}
}

// Test Card struct creation
func TestCard(t *testing.T) {
	tests := []struct {
		name     string
		secret   string
		key      string
		currency string
	}{
		{
			name:     "valid card with USD",
			secret:   "sk_test_123",
			key:      "pk_test_123",
			currency: "usd",
		},
		{
			name:     "valid card with EUR",
			secret:   "sk_test_456",
			key:      "pk_test_456",
			currency: "eur",
		},
		{
			name:     "card with empty values",
			secret:   "",
			key:      "",
			currency: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			card := &Card{
				Secret:   tt.secret,
				Key:      tt.key,
				Currency: tt.currency,
			}

			if card.Secret != tt.secret {
				t.Errorf("Card.Secret = %q, want %q", card.Secret, tt.secret)
			}
			if card.Key != tt.key {
				t.Errorf("Card.Key = %q, want %q", card.Key, tt.key)
			}
			if card.Currency != tt.currency {
				t.Errorf("Card.Currency = %q, want %q", card.Currency, tt.currency)
			}
		})
	}
}

// Test Transaction struct
func TestTransaction(t *testing.T) {
	tests := []struct {
		name                string
		transactionStatusID int
		amount              int
		currency            string
		lastFour            string
		bankReturnCode      string
	}{
		{
			name:                "successful transaction",
			transactionStatusID: 1,
			amount:              1000,
			currency:            "usd",
			lastFour:            "4242",
			bankReturnCode:      "approved",
		},
		{
			name:                "declined transaction",
			transactionStatusID: 2,
			amount:              5000,
			currency:            "eur",
			lastFour:            "0002",
			bankReturnCode:      "declined",
		},
		{
			name:                "zero amount transaction",
			transactionStatusID: 0,
			amount:              0,
			currency:            "",
			lastFour:            "",
			bankReturnCode:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			txn := &Transaction{
				TransactionStatusID: tt.transactionStatusID,
				Amount:              tt.amount,
				Currency:            tt.currency,
				LastFour:            tt.lastFour,
				BankReturnCode:      tt.bankReturnCode,
			}

			if txn.TransactionStatusID != tt.transactionStatusID {
				t.Errorf("Transaction.TransactionStatusID = %d, want %d", txn.TransactionStatusID, tt.transactionStatusID)
			}
			if txn.Amount != tt.amount {
				t.Errorf("Transaction.Amount = %d, want %d", txn.Amount, tt.amount)
			}
			if txn.Currency != tt.currency {
				t.Errorf("Transaction.Currency = %q, want %q", txn.Currency, tt.currency)
			}
			if txn.LastFour != tt.lastFour {
				t.Errorf("Transaction.LastFour = %q, want %q", txn.LastFour, tt.lastFour)
			}
			if txn.BankReturnCode != tt.bankReturnCode {
				t.Errorf("Transaction.BankReturnCode = %q, want %q", txn.BankReturnCode, tt.bankReturnCode)
			}
		})
	}
}

// Test error message consistency
func TestCardErrorMessageConsistency(t *testing.T) {
	// Test that calling the function multiple times with the same code returns the same message
	code := stripe.ErrorCodeCardDeclined
	firstCall := cardErrorMessage(code)

	for i := 0; i < 100; i++ {
		result := cardErrorMessage(code)
		if result != firstCall {
			t.Errorf("cardErrorMessage(%v) returned inconsistent results: first=%q, iteration %d=%q",
				code, firstCall, i, result)
		}
	}
}

// Test all error messages are non-empty
func TestCardErrorMessageNonEmpty(t *testing.T) {
	allCodes := []stripe.ErrorCode{
		stripe.ErrorCodeCardDeclined,
		stripe.ErrorCodeExpiredCard,
		stripe.ErrorCodeIncorrectCVC,
		stripe.ErrorCodeIncorrectZip,
		stripe.ErrorCodeAmountTooLarge,
		stripe.ErrorCodeAmountTooSmall,
		stripe.ErrorCodeBalanceInsufficient,
		stripe.ErrorCodePostalCodeInvalid,
		"unknown_code",
		"",
	}

	for _, code := range allCodes {
		t.Run(string(code), func(t *testing.T) {
			msg := cardErrorMessage(code)
			if msg == "" {
				t.Errorf("cardErrorMessage(%v) returned empty message", code)
			}
		})
	}
}
