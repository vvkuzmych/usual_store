package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"github.com/go-playground/validator/v10"
	"time"
	"usual_store/internal/models"
)

// TokenRepository defines the methods for token-related database operations.
type TokenRepository interface {
	InsertToken(ctx context.Context, token *models.Token, user models.User) error
	GetUserForToken(ctx context.Context, token string) (*models.User, error)
}

type DBModel struct {
	DB *sql.DB
}

// NewDBModel creates a new instance of DBModel.
func NewDBModel(db *sql.DB) *DBModel {
	return &DBModel{DB: db}
}

// InsertToken inserts or updates a token for a user in the database.
func (m *DBModel) InsertToken(ctx context.Context, token *models.Token, user models.User) error {
	validate := validator.New()

	// Validate the struct
	err := validate.Struct(user)
	if err != nil {
		// Loop through validation errors
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Field '%s' failed validation, Condition: '%s'\n", err.Field(), err.ActualTag())
		}
		return err
	}

	// Use $1 for parameterized queries in PostgreSQL
	stmt := `DELETE FROM tokens WHERE user_id = $1`
	_, err = m.DB.ExecContext(ctx, stmt, user.ID)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO tokens 
				(user_id, name, email, token_hash, expiry, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = m.DB.ExecContext(ctx, stmt, user.ID, user.LastName, user.Email, token.Hash, token.Expiry, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

// GetUserForToken retrieves a user based on the provided token.
func (m *DBModel) GetUserForToken(ctx context.Context, token string) (*models.User, error) {
	var user models.User
	tokenHash := sha256.Sum256([]byte(token))

	// PostgreSQL specific query using CURRENT_TIMESTAMP
	query := `SELECT 
		u.id, u.last_name, u.email, u.first_name 
		FROM users u 
		INNER JOIN tokens t ON u.id = t.user_id 
		WHERE t.token_hash = $1 AND t.expiry > $2`

	// Execute the query with $1 and $2 placeholders for the token hash and current time
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:], time.Now()).Scan(
		&user.ID, &user.LastName, &user.Email, &user.FirstName,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
