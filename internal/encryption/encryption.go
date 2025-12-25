package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

type Encryption struct {
	Key []byte
}

// Encrypt encrypts plaintext using AES-GCM (AEAD mode)
func (e *Encryption) Encrypt(plaintext string) (string, error) {
	plainBytes := []byte(plaintext)
	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Use GCM mode (AEAD - Authenticated Encryption with Associated Data)
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a nonce (number used once)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt and authenticate
	ciphertext := aesGCM.Seal(nonce, nonce, plainBytes, nil)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-GCM (AEAD mode)
func (e *Encryption) Decrypt(ciphertext string) (string, error) {
	cipherBytes, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(e.Key)
	if err != nil {
		return "", err
	}

	// Use GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(cipherBytes) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Extract nonce and ciphertext
	nonce, cipherBytes := cipherBytes[:nonceSize], cipherBytes[nonceSize:]

	// Decrypt and verify authentication
	plainBytes, err := aesGCM.Open(nil, nonce, cipherBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plainBytes), nil
}
