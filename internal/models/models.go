package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
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
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
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
	row := m.DB.QueryRowContext(ctx, `SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE email=$1`, email)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
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

// GetOrders retrieves subset orders from the database.
func (m *DBModel) GetOrdersPaginated(isRecurring bool, pageSize, page int) ([]*Order, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	offset := (page - 1) * pageSize

	var orders []*Order

	query := `SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, 
					 o.status_id, o.quantity, o.amount, o.created_at,
					 o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
					 t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
					 t.bank_return_code, c.id, c.first_name, c.last_name, c.email
			  FROM orders o 
			  		LEFT JOIN widgets w ON (o.widget_id = w.id)
					LEFT JOIN transactions t ON (o.transaction_id = t.id)
					LEFT JOIN customers c ON (o.customer_id = c.id)
			  WHERE 
			      w.is_recurring = $1
			  ORDER BY 
			      o.created_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := m.DB.QueryContext(ctx, query, isRecurring, pageSize, offset)
	if err != nil {
		fmt.Println("error getting data from rows:", err)
		return nil, 0, 0, err
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
			fmt.Println("scanning error:", err)
			return nil, 0, 0, err
		}
		orders = append(orders, &order)
	}

	query = `SELECT COUNT(o.id) FROM orders o 
			 LEFT JOIN widgets w ON (o.widget_id = w.id)
			 WHERE w.is_recurring = $1`

	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, query, isRecurring)

	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}
	lastPage := totalRecords / pageSize

	return orders, lastPage, totalRecords, nil
}

// GetAllOrders gets all non-recurring orders from the database.
func (m *DBModel) GetAllOrders(pageSize, page int) ([]*Order, int, int, error) {
	return m.GetOrdersPaginated(false, pageSize, page)
}

// GetAllSubscriptions gets all recurring orders from the database.
func (m *DBModel) GetAllSubscriptions(pageSize, page int) ([]*Order, int, int, error) {
	return m.GetOrdersPaginated(true, pageSize, page)
}

// GetOrderByID gets order by id
func (m *DBModel) GetOrderByID(id int) (Order, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var order Order

	query := `SELECT o.id, o.widget_id, o.transaction_id, o.customer_id, 
					 o.status_id, o.quantity, o.amount, o.created_at,
					 o.updated_at, w.id, w.name, t.id, t.amount, t.currency,
					 t.last_four, t.expiry_month, t.expiry_year, t.payment_intent,
					 t.bank_return_code, c.id, c.first_name, c.last_name, c.email
			  FROM orders o 
			  		LEFT JOIN widgets w ON (o.widget_id = w.id)
					LEFT JOIN transactions t ON (o.transaction_id = t.id)
					LEFT JOIN customers c ON (o.customer_id = c.id)
			  WHERE 
			      o.id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
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
		fmt.Println("scanning error:", err)
		return order, err
	}

	return order, nil
}

// UpdateOrderStatus updates order status
func (m *DBModel) UpdateOrderStatus(id, statusID int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "UPDATE orders SET status_id = $1 WHERE id = $2"
	_, err := m.DB.ExecContext(ctx, stmt, statusID, id)
	if err != nil {
		return err
	}
	return nil
}

// GetAllUsers gets all users
func (m *DBModel) GetAllUsers() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var users []*User
	query := `SELECT id, first_name, last_name, email, created_at, updated_at FROM users ORDER BY last_name, first_name DESC`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	return users, nil
}

// GetUserByID gets user by id
func (m *DBModel) GetUserByID(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User

	query := `SELECT id, first_name, last_name, email, password, role, created_at, updated_at FROM users WHERE id = $1`
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, fmt.Errorf("no user found with ID %d", id)
		}
		return user, err
	}
	return user, nil
}

// EditUser update user data
func (m *DBModel) EditUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `UPDATE users SET first_name = $1, last_name = $2, email = $3, updated_at = $4 WHERE id = $5`
	_, err := m.DB.ExecContext(ctx, stmt, user.FirstName, user.LastName, user.Email, time.Now(), user.ID)

	if err != nil {
		return err
	}
	return nil
}

// AddUser add user
func (m *DBModel) AddUser(user User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Default role to 'user' if not specified
	role := user.Role
	if role == "" {
		role = "user"
	}

	stmt := `INSERT INTO users (first_name, last_name, email, password, role, created_at, updated_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(ctx, stmt, user.FirstName, user.LastName, user.Email, hash, role, time.Now(), time.Now())

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser remove user
func (m *DBModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM users WHERE id = $1`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = `DELETE FROM tokens WHERE user_id = $1`
	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllUsersPaginated returns paginated list of all users
func (m *DBModel) GetAllUsersPaginated(offset, limit int) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT id, first_name, last_name, email, role, created_at, updated_at 
			  FROM users 
			  ORDER BY 
			    CASE role
			      WHEN 'super_admin' THEN 1
			      WHEN 'admin' THEN 2
			      WHEN 'supporter' THEN 3
			      ELSE 4
			    END, id
			  LIMIT $1 OFFSET $2`

	rows, err := m.DB.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// GetUserCount returns the total count of users
func (m *DBModel) GetUserCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// GetSuperAdminCount returns the count of super_admin users
func (m *DBModel) GetSuperAdminCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	err := m.DB.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE role = 'super_admin'`).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
