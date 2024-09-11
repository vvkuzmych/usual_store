package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base32"
	"log"
	"time"
)

const (
	ScopeAuthentication = "authentication"
)

// Token is the type of authentication token
type Token struct {
	PlainText string    `json:"token"`
	UserId    int64     `json:"-"`
	Hash      []byte    `json:"-"`
	Expiry    time.Time `json:"expiry"`
	Scope     string    `json:"-"`
}

// GenerateToken generates token that lasts ttl and returns it
func GenerateToken(userId int, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserId: int64(userId),
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.PlainText = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	hash := sha256.Sum256([]byte(token.PlainText))
	token.Hash = hash[:]
	return token, nil
}

func (m *DBModel) InsertToken(token *Token, user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `DELETE FROM tokens WHERE user_id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, user.ID)
	if err != nil {
		return err
	}

	stmt = `INSERT INTO tokens 
				(user_id, name, email, token_hash, created_at, updated_at)
				values(?, ?, ?, ?, ?, ?)
				`
	_, err = m.DB.ExecContext(ctx, stmt, user.ID, user.LastName, user.Email, token.Hash, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) GetUserForToken(token string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	var user User
	tokenHash := sha256.Sum256([]byte(token))

	query := `SELECT 
    u.id, u.last_name, u.email, u.first_name 
    FROM users u 
    INNER JOIN tokens t ON u.user_id = t.user_id 
    WHERE t.hash = ?`
	err := m.DB.QueryRowContext(ctx, query, tokenHash[:]).Scan(
		&user.ID, &user.LastName, &user.Email, &user.FirstName,
	)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &user, nil
}
