package service

import (
	"database/sql"
	"encoding/json"
	"time"
)

// PostgresSessionService is a custom PostgreSQL session store implementation.
type PostgresSessionService struct {
	db *sql.DB
}

// NewPostgresSessionService creates a new instance of PostgresSessionService.
func NewPostgresSessionService(db *sql.DB) *PostgresSessionService {
	return &PostgresSessionService{db: db}
}

// SessionService defines the methods that the session service must implement.
type SessionService interface {
	Find(token string) (b []byte, found bool, err error)
	Commit(token string, b []byte, expiry time.Time) (err error)
	Get(sessionID string) (map[string]interface{}, error)
	Set(sessionID string, sessionData map[string]interface{}, expiry time.Time) error
	Delete(sessionID string) error
}

// Find retrieves a session from the PostgreSQL database.
func (s *PostgresSessionService) Find(token string) (b []byte, found bool, err error) {
	var data []byte
	var expiry time.Time

	err = s.db.QueryRow("SELECT data, expiry FROM sessions WHERE token = $1", token).Scan(&data, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false, nil
		}
		return nil, false, err
	}

	if expiry.Before(time.Now()) {
		return nil, false, nil
	}

	return data, true, nil
}

// Commit stores or updates a session in the PostgreSQL database.
func (s *PostgresSessionService) Commit(token string, b []byte, expiry time.Time) (err error) {
	_, err = s.db.Exec(
		"INSERT INTO sessions (token, data, expiry) VALUES ($1, $2, $3) "+
			"ON CONFLICT (token) DO UPDATE SET data = $2, expiry = $3",
		token, b, expiry,
	)
	return err
}

// Get retrieves a session from the PostgreSQL database.
func (s *PostgresSessionService) Get(sessionID string) (map[string]interface{}, error) {
	var data []byte
	var expiry time.Time

	err := s.db.QueryRow("SELECT data, expiry FROM sessions WHERE token = $1", sessionID).Scan(&data, &expiry)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if expiry.Before(time.Now()) {
		return nil, nil
	}

	var sessionData map[string]interface{}
	err = json.Unmarshal(data, &sessionData)
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

// Set saves a session to the PostgreSQL database.
func (s *PostgresSessionService) Set(sessionID string, sessionData map[string]interface{}, expiry time.Time) error {
	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO sessions (token, data, expiry) VALUES ($1, $2, $3) "+
			"ON CONFLICT (token) DO UPDATE SET data = $2, expiry = $3",
		sessionID, data, expiry,
	)
	return err
}

// Delete removes a session from the PostgreSQL database.
func (s *PostgresSessionService) Delete(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE token = $1", sessionID)
	return err
}
