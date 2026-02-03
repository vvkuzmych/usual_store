package models

import (
	"context"
	"fmt"
	"time"
)

// Customer represents a customer in the system.
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// InsertCustomer adds a new customer to the database.
func (m *DBModel) InsertCustomer(customer Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert the customer into the database
	const insertQuery = `
		INSERT INTO customers (first_name, last_name, email)
		VALUES ($1, $2, $3)
	`

	_, err := m.DB.ExecContext(ctx, insertQuery, customer.FirstName, customer.LastName, customer.Email)
	if err != nil {
		return fmt.Errorf("failed to insert customer: %w", err)
	}

	return nil
}
