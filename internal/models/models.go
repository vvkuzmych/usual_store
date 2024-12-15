package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"strings"
	"time"
)

// DBModel represents a model with a database connection.
type DBModel struct {
	DB *sql.DB
}

// Models is the type for all models
type Models struct {
	DB DBModel
}

// NewModels returns a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for all widgets
type Widget struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	IsRecurring    bool      `json:"is_recurring"`
	PlanID         string    `json:"plan_id"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for all orders
type Order struct {
	ID            int         `json:"id"`
	WidgetID      int         `json:"widget_id"`
	TransactionID int         `json:"transaction_id"`
	CustomerID    int         `json:"customer_id"`
	StatusID      int         `json:"status_id"`
	Quantity      int         `json:"quantity"`
	Amount        int         `json:"amount"`
	Widget        Widget      `json:"widget"`
	Transaction   Transaction `json:"transaction"`
	Customer      Customer    `json:"customer"`
	CreatedAt     time.Time   `json:"-"`
	UpdatedAt     time.Time   `json:"-"`
}

// Status is the type for statuses
type Status struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for transaction statuses
type TransactionStatus struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for transactions
type Transaction struct {
	ID                  int       `json:"id"`
	Amount              int       `json:"amount"`
	Currency            string    `json:"currency"`
	LastFour            string    `json:"last_four"`
	ExpiryMonth         int       `json:"expiry_month"`
	ExpiryYear          int       `json:"expiry_year"`
	BankReturnCode      string    `json:"bank_return_code"`
	TransactionStatusID int       `json:"transaction_status_id"`
	PaymentIntent       string    `json:"payment_intent"`
	PaymentMethod       string    `json:"payment_method"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// User is the type for users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Customer is the type for customers
type Customer struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// GetWidget gets widget by id
func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `SELECT id, name, description, inventory_level, price, image, is_recurring, plan_id, created_at, updated_at FROM widgets WHERE id=$1`
	row := m.DB.QueryRowContext(ctx, stmt, id)

	var widget Widget
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.IsRecurring,
		&widget.PlanID,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return widget, fmt.Errorf("no widget found with id %d", id)
		}
		return widget, err
	}

	return widget, nil
}

// InsertTransaction insert a new txn and returns new id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Insert the transaction into the database
	stmt := `INSERT INTO transactions 
				(amount, currency, last_four, bank_return_code, 
				 expiry_month, expiry_year, payment_intent, payment_method, 
				 transaction_status_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	_, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.BankReturnCode,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentIntent,
		txn.PaymentMethod,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return 0, err
	}

	// Retrieve the transaction ID using the paymentIntent
	id, err := m.GetTransactionIDByPaymentIntent(ctx, txn.PaymentIntent)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// GetTransactionIDByPaymentIntent retrieves the transaction ID based on the payment_intent.
func (m *DBModel) GetTransactionIDByPaymentIntent(ctx context.Context, paymentIntent string) (int, error) {
	stmt := `SELECT id FROM transactions WHERE payment_intent = $1;`

	var id int
	err := m.DB.QueryRowContext(ctx, stmt, paymentIntent).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// InsertOrder insert a new order and returns new id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Check if widget_id exists
	err := m.CheckWidgetExistence(ctx, order.WidgetID)
	if err != nil {
		return 0, fmt.Errorf("widget does not exist: %w", err)
	}

	// Insert the order into the database
	stmt := `
        INSERT INTO orders 
        (widget_id, transaction_id, status_id, quantity, customer_id, amount)
        VALUES($1, $2, $3, $4, $5, $6)
    `
	_, err = m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.Quantity,
		order.CustomerID,
		order.Amount,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert order: %w", err)
	}

	// Retrieve the latest inserted ID
	var orderID int
	query := `
        SELECT id FROM orders
        WHERE widget_id = $1 AND transaction_id = $2 AND customer_id = $3
        ORDER BY created_at DESC
        LIMIT 1
    `
	err = m.DB.QueryRowContext(ctx, query, order.WidgetID, order.TransactionID, order.CustomerID).Scan(&orderID)
	if err != nil {
		return 0, fmt.Errorf("failed to retrieve order ID: %w", err)
	}

	return orderID, nil
}

// CheckWidgetExistence checks if a widget exists based on its ID.
func (m *DBModel) CheckWidgetExistence(ctx context.Context, widgetID int) error {
	var widgetExists bool
	err := m.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM widgets WHERE id = $1)", widgetID).Scan(&widgetExists)
	if err != nil {
		return fmt.Errorf("could not check widget existence: %v", err)
	}
	if !widgetExists {
		return fmt.Errorf("widget_id %d does not exist", widgetID)
	}
	return nil
}

// InsertCustomer inserts a customer record into the database.
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
		logQueryError("InsertCustomer", insertQuery, customer, err)
		return fmt.Errorf("failed to insert customer: %w", err)
	}

	return nil
}

// GetLastInsertedCustomerID retrieves the last inserted customer ID.
func (m *DBModel) GetLastInsertedCustomerID() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Query to fetch the last inserted customer ID
	const getIDQuery = `
		SELECT id
		FROM customers
		ORDER BY id DESC
		LIMIT 1
	`

	var customerID int
	err := m.DB.QueryRowContext(ctx, getIDQuery).Scan(&customerID)
	if err != nil {
		logQueryError("GetLastInsertedCustomerID", getIDQuery, nil, err)
		return 0, fmt.Errorf("failed to get last inserted customer ID: %w", err)
	}

	return customerID, nil
}

// logQueryError logs detailed information about query execution failures.
func logQueryError(operation, query string, params interface{}, err error) {
	log.Printf("[%s] Query failed.\nQuery: %s\nParams: %+v\nError: %v", operation, query, params, err)
}

// CheckCustomerExistence checks if a customer exists based on their email.
func (m *DBModel) CheckCustomerExistence(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := m.DB.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM customers WHERE email = $1)", email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("could not check customer existence: %v", err)
	}
	return exists, nil
}

// GetUserByEmail gets user by email address
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var user User
	// Update query to use $1 for parameterized query in PostgreSQL
	row := m.DB.QueryRowContext(ctx, `SELECT id, first_name, last_name, email, password, created_at, updated_at FROM users WHERE email=$1`, email)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no user found with email %s", email)
		}
		return user, err
	}

	return user, nil
}

func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email=$1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return id, fmt.Errorf("no user found with email %s", email)
		}
		return id, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return 0, fmt.Errorf("invalid password")
	} else if err != nil {
		return 0, err
	}
	return id, nil
}

// UpdatePasswordForUser it updates password for user
func (m *DBModel) UpdatePasswordForUser(user User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
	UPDATE users
	SET password = $1
	WHERE id = $2
	`

	_, err := m.DB.ExecContext(ctx, stmt, hash, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) GetAllOrders() ([]*Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var orders []*Order

	query := `SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, 
					 o.status_id, o.quantity, o.amount, o.created_at,
					 o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
					 t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
					 t.bank_return_code, c.id, c.first_name, c.last_name, c.email
			  FROM orders o 
			  		left join widgets w on (o.widget_id = w.id)
					left join transactions t on (o.transaction_id = t.id)
					left join customers c on (o.customer_id = c.id)
			  WHERE 
			      w.is_recurring = false
			  ORDER BY 
			      o.created_at desc`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		fmt.Println("error getting data from rows: ", err)

		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var order Order
		err = rows.Scan(
			&order.ID,
			&order.WidgetID,
			&order.TransactionID,
			&order.CustomerID,
			&order.StatusID,
			&order.Quantity,
			&order.Amount,
			&order.CreatedAt,
			&order.UpdatedAt,
			&order.Widget.ID,
			&order.Widget.Name,
			&order.Transaction.ID,
			&order.Transaction.Amount,
			&order.Transaction.Currency,
			&order.Transaction.LastFour,
			&order.Transaction.ExpiryMonth,
			&order.Transaction.ExpiryYear,
			&order.Transaction.PaymentIntent,
			&order.Transaction.BankReturnCode,
			&order.Customer.ID,
			&order.Customer.FirstName,
			&order.Customer.LastName,
			&order.Customer.Email,
		)
		if err != nil {
			fmt.Println("scanning error")

			return nil, err
		}
		orders = append(orders, &order)

	}
	return orders, nil
}
