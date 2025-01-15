package models

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCheckCustomerExistence(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	assert.NoError(t, err, "unexpected error creating mock DB")
	defer db.Close()

	// Create the Models instance
	models := NewModels(db)

	// Define test cases
	sqlQuery := "SELECT EXISTS\\(SELECT 1 FROM customers WHERE email = \\$1\\)"
	tests := []struct {
		name       string
		email      string
		mockQuery  func()
		expected   bool
		shouldFail bool
		err        error
	}{
		{
			name:  "Customer exists",
			email: "customer@example.com",
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs("customer@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
			expected:   true,
			shouldFail: false,
		},
		{
			name:  "Customer does not exist",
			email: "nonexistent@example.com",
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs("nonexistent@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"does not exists"}).AddRow(false))
			},
			expected:   false,
			shouldFail: false,
		},
		{
			name:  "Database error",
			email: "error@example.com",
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs("error@example.com").
					WillReturnError(errors.New("database error"))
			},
			expected:   false,
			shouldFail: true,
			err:        errors.New("database error"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockQuery()

			// Execute the function
			result, err := models.DB.CheckCustomerExistence(context.Background(), tc.email)

			// Assert results
			if tc.shouldFail {
				assert.Contains(t, err.Error(), tc.err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, result)

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
		})
	}
}
