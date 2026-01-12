package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"testing"
	"time"
	"usual_store/internal/models"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBModel_InsertToken(t *testing.T) {
	tests := []struct {
		name      string
		token     *models.Token
		user      models.User
		mockSetup func(mock sqlmock.Sqlmock, token *models.Token, user models.User)
		validate  func(t *testing.T, err error, mockErr error)
	}{
		{
			name: "successful token insertion",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
				Password:  "hashed-password",
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// Expect DELETE query for old tokens
				mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// Expect INSERT query for new token
				mock.ExpectExec("INSERT INTO tokens").
					WithArgs(
						user.ID,
						user.LastName,
						user.Email,
						token.Hash,
						sqlmock.AnyArg(), // expiry
						sqlmock.AnyArg(), // created_at
						sqlmock.AnyArg(), // updated_at
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mockErr error) {
				assert.NoError(t, err)
				assert.NoError(t, mockErr)
			},
		},
		{
			name: "no tokens to delete before insert",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        2,
				FirstName: "Jane",
				LastName:  "Smith",
				Email:     "jane@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// DELETE returns 0 rows affected (no existing tokens)
				mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))

				// INSERT new token
				mock.ExpectExec("INSERT INTO tokens").
					WithArgs(
						user.ID,
						user.LastName,
						user.Email,
						token.Hash,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			validate: func(t *testing.T, err error, mockErr error) {
				assert.NoError(t, err)
				assert.NoError(t, mockErr)
			},
		},
		{
			name: "database error on delete",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// DELETE fails
				mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
					WithArgs(user.ID).
					WillReturnError(errors.New("database connection error"))
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "database connection error")
			},
		},
		{
			name: "database error on insert",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// DELETE succeeds
				mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))

				// INSERT fails
				mock.ExpectExec("INSERT INTO tokens").
					WithArgs(
						user.ID,
						user.LastName,
						user.Email,
						token.Hash,
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
						sqlmock.AnyArg(),
					).
					WillReturnError(errors.New("insert failed"))
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "insert failed")
			},
		},
		{
			name: "validation error - missing last name",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "", // Empty last name should fail validation
				Email:     "john@example.com",
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// No mock expectations - validation should fail before DB calls
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "LastName")
			},
		},
		{
			name: "validation error - missing email",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "", // Empty email should fail validation
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// No mock expectations - validation should fail before DB calls
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Email")
			},
		},
		{
			name: "validation error - invalid email format",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
				Email:     "not-an-email", // Invalid email format
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// No mock expectations - validation should fail before DB calls
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "Email")
			},
		},
		{
			name: "validation error - multiple validation failures",
			token: &models.Token{
				Hash:   []byte("test-hash"),
				Expiry: time.Now().Add(24 * time.Hour),
			},
			user: models.User{
				ID:        1,
				FirstName: "John",
				LastName:  "",              // Missing
				Email:     "invalid-email", // Invalid format
			},
			mockSetup: func(mock sqlmock.Sqlmock, token *models.Token, user models.User) {
				// No mock expectations - validation should fail before DB calls
			},
			validate: func(t *testing.T, err error, mockErr error) {
				require.Error(t, err)
				// Should contain at least one field error
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			defer db.Close()

			// Setup mock expectations
			tt.mockSetup(mock, tt.token, tt.user)

			// Create repository
			repo := NewDBModel(db)

			// Execute test
			ctx := context.Background()
			err = repo.InsertToken(ctx, tt.token, tt.user)

			// Check mock expectations
			mockErr := mock.ExpectationsWereMet()

			// Validate results
			tt.validate(t, err, mockErr)
		})
	}
}

func TestDBModel_GetUserForToken(t *testing.T) {
	tests := []struct {
		name      string
		token     string
		mockSetup func(mock sqlmock.Sqlmock, tokenHash [32]byte)
		validate  func(t *testing.T, user *models.User, err error)
	}{
		{
			name:  "successful user retrieval with valid token",
			token: "valid-token-12345",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				rows := sqlmock.NewRows([]string{"id", "last_name", "email", "first_name"}).
					AddRow(1, "Doe", "john@example.com", "John")

				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, 1, user.ID)
				assert.Equal(t, "John", user.FirstName)
				assert.Equal(t, "Doe", user.LastName)
				assert.Equal(t, "john@example.com", user.Email)
			},
		},
		{
			name:  "user retrieval with different user data",
			token: "another-valid-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				rows := sqlmock.NewRows([]string{"id", "last_name", "email", "first_name"}).
					AddRow(42, "Smith", "jane.smith@example.com", "Jane")

				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, 42, user.ID)
				assert.Equal(t, "Jane", user.FirstName)
				assert.Equal(t, "Smith", user.LastName)
				assert.Equal(t, "jane.smith@example.com", user.Email)
			},
		},
		{
			name:  "token not found in database",
			token: "non-existent-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.True(t, errors.Is(err, sql.ErrNoRows), "error should be sql.ErrNoRows")
				assert.Nil(t, user)
			},
		},
		{
			name:  "expired token returns no rows",
			token: "expired-token-xyz",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.True(t, errors.Is(err, sql.ErrNoRows), "error should be sql.ErrNoRows")
				assert.Nil(t, user)
			},
		},
		{
			name:  "database connection error",
			token: "some-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				dbErr := errors.New("connection refused")
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnError(dbErr)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "connection refused")
				assert.Nil(t, user)
			},
		},
		{
			name:  "database query timeout",
			token: "timeout-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				timeoutErr := errors.New("query timeout exceeded")
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnError(timeoutErr)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "timeout")
				assert.Nil(t, user)
			},
		},
		{
			name:  "scan error - insufficient columns",
			token: "malformed-result-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				rows := sqlmock.NewRows([]string{"id", "last_name"}).
					AddRow(1, "Doe")

				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name:  "scan error - type mismatch",
			token: "type-mismatch-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				rows := sqlmock.NewRows([]string{"id", "last_name", "email", "first_name"}).
					AddRow("not-an-int", "Doe", "john@example.com", "John")

				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.Nil(t, user)
			},
		},
		{
			name:  "empty token string",
			token: "",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnError(sql.ErrNoRows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				require.Error(t, err)
				assert.True(t, errors.Is(err, sql.ErrNoRows), "error should be sql.ErrNoRows")
				assert.Nil(t, user)
			},
		},
		{
			name:  "user with empty first name",
			token: "no-firstname-token",
			mockSetup: func(mock sqlmock.Sqlmock, tokenHash [32]byte) {
				rows := sqlmock.NewRows([]string{"id", "last_name", "email", "first_name"}).
					AddRow(5, "Anonymous", "anonymous@example.com", "")

				mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
					WithArgs(tokenHash[:], sqlmock.AnyArg()).
					WillReturnRows(rows)
			},
			validate: func(t *testing.T, user *models.User, err error) {
				assert.NoError(t, err)
				require.NotNil(t, user)
				assert.Equal(t, 5, user.ID)
				assert.Equal(t, "", user.FirstName)
				assert.Equal(t, "Anonymous", user.LastName)
				assert.Equal(t, "anonymous@example.com", user.Email)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock database
			db, mock, err := sqlmock.New()
			require.NoError(t, err, "failed to create mock database")
			defer db.Close()

			// Calculate token hash
			tokenHash := sha256.Sum256([]byte(tt.token))

			// Setup mock expectations
			tt.mockSetup(mock, tokenHash)

			// Create repository
			repo := NewDBModel(db)
			require.NotNil(t, repo)

			// Execute test
			ctx := context.Background()
			user, err := repo.GetUserForToken(ctx, tt.token)

			// Validate results
			tt.validate(t, user, err)

			// Ensure all mock expectations were met
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestDBModel_InsertToken_ContextCancellation(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	token := &models.Token{
		Hash:   []byte("test-hash"),
		Expiry: time.Now().Add(24 * time.Hour),
	}
	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	// Expect DELETE to fail due to context cancellation
	mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
		WithArgs(user.ID).
		WillReturnError(context.Canceled)

	repo := NewDBModel(db)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = repo.InsertToken(ctx, token, user)

	// Should return context cancellation error
	assert.Error(t, err)
}

func TestDBModel_GetUserForToken_ContextCancellation(t *testing.T) {
	// Create mock database
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	token := "test-token"
	tokenHash := sha256.Sum256([]byte(token))

	// Expect query to fail due to context cancellation
	mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
		WithArgs(tokenHash[:], sqlmock.AnyArg()).
		WillReturnError(context.Canceled)

	repo := NewDBModel(db)

	// Create cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	user, err := repo.GetUserForToken(ctx, token)

	// Should return context cancellation error
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestNewDBModel(t *testing.T) {
	db, _, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewDBModel(db)

	assert.NotNil(t, repo)
	assert.NotNil(t, repo.DB)
	assert.Equal(t, db, repo.DB)
}

// Benchmark tests
func BenchmarkInsertToken(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	token := &models.Token{
		Hash:   []byte("test-hash"),
		Expiry: time.Now().Add(24 * time.Hour),
	}
	user := models.User{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@example.com",
	}

	// Setup expectations for N iterations
	for i := 0; i < b.N; i++ {
		mock.ExpectExec("DELETE FROM tokens WHERE user_id = \\$1").
			WithArgs(user.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		mock.ExpectExec("INSERT INTO tokens").
			WithArgs(
				user.ID,
				user.LastName,
				user.Email,
				token.Hash,
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
	}

	repo := NewDBModel(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = repo.InsertToken(ctx, token, user)
	}
}

func BenchmarkGetUserForToken(b *testing.B) {
	db, mock, err := sqlmock.New()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	token := "test-token"
	tokenHash := sha256.Sum256([]byte(token))

	// Setup expectations for N iterations
	for i := 0; i < b.N; i++ {
		rows := sqlmock.NewRows([]string{"id", "last_name", "email", "first_name"}).
			AddRow(1, "Doe", "john@example.com", "John")

		mock.ExpectQuery("SELECT (.+) FROM users u INNER JOIN tokens t").
			WithArgs(tokenHash[:], sqlmock.AnyArg()).
			WillReturnRows(rows)
	}

	repo := NewDBModel(db)
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.GetUserForToken(ctx, token)
	}
}
