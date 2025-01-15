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
			expected: true,
		},
		{
			name:  "Customer does not exist",
			email: "nonexistent@example.com",
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs("nonexistent@example.com").
					WillReturnRows(sqlmock.NewRows([]string{"does not exists"}).AddRow(false))
			},
		},
		{
			name:  "Database error",
			email: "error@example.com",
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs("error@example.com").
					WillReturnError(errors.New("database error"))
			},
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

func TestCheckWidgetExistence(t *testing.T) {
	// Mock database setup
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlQuery := "SELECT EXISTS\\(SELECT 1 FROM widgets WHERE id = \\$1\\)"
	tests := []struct {
		name       string
		widgetID   int
		mockQuery  func()
		expected   error
		shouldFail bool
	}{
		{
			name:     "Widget exists",
			widgetID: 1,
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs(1).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
			},
		},
		{
			name:     "Widget does not exist",
			widgetID: 2,
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs(2).
					WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
			},
			expected:   errors.New("widget_id 2 does not exist"),
			shouldFail: true,
		},
		{
			name:     "Database error",
			widgetID: 3,
			mockQuery: func() {
				mock.ExpectQuery(sqlQuery).
					WithArgs(3).
					WillReturnError(errors.New("database error"))
			},
			expected:   errors.New("could not check widget existence: database error"),
			shouldFail: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Apply test-specific mock behavior
			tc.mockQuery()

			// Create the DBModel instance
			model := DBModel{DB: db}

			// Execute the function
			err = model.CheckWidgetExistence(context.Background(), tc.widgetID)

			// Assert results
			if tc.shouldFail {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expected.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, err)
			}

			// Ensure all expectations were met
			assert.NoError(t, mock.ExpectationsWereMet(), "there were unfulfilled expectations")
		})
	}
}
