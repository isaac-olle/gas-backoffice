package infrastructure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/isaac-olle/gas-backoffice/internal/shared/persistance"
)

// AESEncrypter is an AES-256-GCM implementation of persistance.Encrypter.
// GCM is authenticated encryption: it both hides the plaintext and detects
// tampering. A fresh random nonce is generated per Encrypt and prepended to the
// ciphertext; the whole thing is base64-encoded for storage as TEXT.
type AESEncrypter struct {
	gcm cipher.AEAD
}

var _ persistance.Encrypter = (*AESEncrypter)(nil)

// NewAESEncrypter builds the encrypter from a base64-encoded 32-byte key
// (AES-256). It fails fast if the key is missing or the wrong size.
func NewAESEncrypter(keyBase64 string) (*AESEncrypter, error) {
	key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil {
		return nil, fmt.Errorf("decode encryption key: %w", err)
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes (AES-256), got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm: %w", err)
	}

	return &AESEncrypter{gcm: gcm}, nil
}

// Seal xifra i retorna nonce||ciphertext cru, sense codificar. Qui crida tria
// com el transporta.
func (e *AESEncrypter) Seal(plaintext []byte) ([]byte, error) {
	nonce := make([]byte, e.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	// Seal appends the ciphertext to nonce, so the result is nonce||ciphertext.
	return e.gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func (e *AESEncrypter) Open(sealed []byte) ([]byte, error) {
	nonceSize := e.gcm.NonceSize()
	if len(sealed) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, body := sealed[:nonceSize], sealed[nonceSize:]
	plaintext, err := e.gcm.Open(nil, nonce, body, nil)
	if err != nil {
		return nil, fmt.Errorf("decrypt: %w", err)
	}
	return plaintext, nil
}

// Encrypt és Seal amb base64 estàndard, per desar en una columna TEXT.
func (e *AESEncrypter) Encrypt(plaintext string) (string, error) {
	sealed, err := e.Seal([]byte(plaintext))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(sealed), nil
}

func (e *AESEncrypter) Decrypt(ciphertext string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode ciphertext: %w", err)
	}

	plaintext, err := e.Open(data)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
