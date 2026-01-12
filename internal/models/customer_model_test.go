package models

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestDBModel_InsertCustomer(t *testing.T) {
	tests := []struct {
		name      string
		customer  Customer
		mockSetup func(mock sqlmock.Sqlmock, customer Customer)
		validate  func(t *testing.T, err error, mock sqlmock.Sqlmock)
	}{
		{
			name: "successful customer insertion",
			customer: Customer{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "successful insertion with different customer data",
			customer: Customer{
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane.smith@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(2, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with empty first name",
			customer: Customer{
				FirstName: "",
				LastName:  "Doe",
				Email:     "empty@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with long names",
			customer: Customer{
				FirstName: "Bartholomew",
				LastName:  "Montgomery-Fitzwilliam",
				Email:     "bartholomew.montgomery@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "database error - duplicate email",
			customer: Customer{
				FirstName: "John",
				LastName:  "Doe",
				Email:     "duplicate@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnError(errors.New("duplicate key value violates unique constraint"))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to insert customer")
				require.Contains(t, err.Error(), "duplicate key")
			},
		},
		{
			name: "database error - connection failure",
			customer: Customer{
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnError(errors.New("connection refused"))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to insert customer")
				require.Contains(t, err.Error(), "connection refused")
			},
		},
		{
			name: "database error - timeout",
			customer: Customer{
				FirstName: "Timeout",
				LastName:  "Test",
				Email:     "timeout@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnError(context.DeadlineExceeded)
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to insert customer")
			},
		},
		{
			name: "customer with special characters in email",
			customer: Customer{
				FirstName: "Special",
				LastName:  "Char",
				Email:     "special+test@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with unicode characters",
			customer: Customer{
				FirstName: "François",
				LastName:  "José",
				Email:     "francois.jose@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "database error - invalid email format caught by database",
			customer: Customer{
				FirstName: "Invalid",
				LastName:  "Email",
				Email:     "not-an-email",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnError(errors.New("invalid email format"))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.Error(t, err)
				require.Contains(t, err.Error(), "failed to insert customer")
			},
		},
		{
			name: "database error - transaction error",
			customer: Customer{
				FirstName: "Transaction",
				LastName:  "Error",
				Email:     "transaction@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnError(sql.ErrTxDone)
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.Error(t, err)
				require.True(t, errors.Is(err, sql.ErrTxDone) || errors.Unwrap(err) == sql.ErrTxDone)
			},
		},
		{
			name: "customer with minimal data",
			customer: Customer{
				FirstName: "A",
				LastName:  "B",
				Email:     "a@b.c",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.customer)

			dbModel := &DBModel{DB: db}
			err = dbModel.InsertCustomer(tt.customer)

			tt.validate(t, err, mock)
		})
	}
}

func TestDBModel_InsertCustomer_ContextTimeout(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	customer := Customer{
		FirstName: "Context",
		LastName:  "Timeout",
		Email:     "context@example.com",
	}

	// Simulate a timeout by expecting an error
	mock.ExpectExec("INSERT INTO customers").
		WithArgs(customer.FirstName, customer.LastName, customer.Email).
		WillDelayFor(5 * time.Second).
		WillReturnError(context.DeadlineExceeded)

	dbModel := &DBModel{DB: db}
	err = dbModel.InsertCustomer(customer)

	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to insert customer")
}

func TestCustomer_StructValidation(t *testing.T) {
	tests := []struct {
		name     string
		customer Customer
		validate func(t *testing.T, customer Customer)
	}{
		{
			name: "complete customer data",
			customer: Customer{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			validate: func(t *testing.T, customer Customer) {
				require.Equal(t, 1, customer.ID)
				require.Equal(t, "John", customer.FirstName)
				require.Equal(t, "Doe", customer.LastName)
				require.Equal(t, "john@example.com", customer.Email)
				require.False(t, customer.CreatedAt.IsZero())
				require.False(t, customer.UpdatedAt.IsZero())
			},
		},
		{
			name: "customer with zero ID",
			customer: Customer{
				ID:        0,
				FirstName: "Test",
				LastName:  "User",
				Email:     "test@example.com",
			},
			validate: func(t *testing.T, customer Customer) {
				require.Equal(t, 0, customer.ID)
				require.NotEmpty(t, customer.FirstName)
				require.NotEmpty(t, customer.LastName)
				require.NotEmpty(t, customer.Email)
			},
		},
		{
			name: "customer with zero time values",
			customer: Customer{
				ID:        5,
				FirstName: "Zero",
				LastName:  "Time",
				Email:     "zero@example.com",
			},
			validate: func(t *testing.T, customer Customer) {
				require.True(t, customer.CreatedAt.IsZero())
				require.True(t, customer.UpdatedAt.IsZero())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.customer)
		})
	}
}

// Benchmark tests
func BenchmarkInsertCustomer(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	customer := Customer{
		FirstName: "Bench",
		LastName:  "Mark",
		Email:     "benchmark@example.com",
	}

	for i := 0; i < b.N; i++ {
		mock.ExpectExec("INSERT INTO customers").
			WithArgs(customer.FirstName, customer.LastName, customer.Email).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	dbModel := &DBModel{DB: db}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dbModel.InsertCustomer(customer)
	}
}

func TestCustomer_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		customer  Customer
		mockSetup func(mock sqlmock.Sqlmock, customer Customer)
		validate  func(t *testing.T, err error, mock sqlmock.Sqlmock)
	}{
		{
			name: "customer with very long email",
			customer: Customer{
				FirstName: "Long",
				LastName:  "Email",
				Email:     "very.long.email.address.that.exceeds.normal.length@subdomain.example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with numbers in name",
			customer: Customer{
				FirstName: "John123",
				LastName:  "Doe456",
				Email:     "john123@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with hyphenated last name",
			customer: Customer{
				FirstName: "Mary",
				LastName:  "Smith-Johnson",
				Email:     "mary.smith-johnson@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
		{
			name: "customer with apostrophe in name",
			customer: Customer{
				FirstName: "O'Brien",
				LastName:  "D'Angelo",
				Email:     "obrien@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, customer Customer) {
				mock.ExpectExec("INSERT INTO customers").
					WithArgs(customer.FirstName, customer.LastName, customer.Email).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mock sqlmock.Sqlmock) {
				require.NoError(t, err)
				require.NoError(t, mock.ExpectationsWereMet())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			tt.mockSetup(mock, tt.customer)

			dbModel := &DBModel{DB: db}
			err = dbModel.InsertCustomer(tt.customer)

			tt.validate(t, err, mock)
		})
	}
}
