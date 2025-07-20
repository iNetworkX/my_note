package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltSize   = 32
	nonceSize  = 12
	keySize    = 32
	iterations = 100000
)

type CryptoManager struct {
	salt []byte
}

func NewCryptoManager() *CryptoManager {
	return &CryptoManager{}
}

func (c *CryptoManager) GenerateSalt() error {
	c.salt = make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, c.salt); err != nil {
		return err
	}
	return nil
}

func (c *CryptoManager) SetSalt(salt []byte) {
	c.salt = salt
}

func (c *CryptoManager) GetSalt() []byte {
	return c.salt
}

func (c *CryptoManager) DeriveKey(password string) []byte {
	return pbkdf2.Key([]byte(password), c.salt, iterations, keySize, sha256.New)
}

func (c *CryptoManager) Encrypt(plaintext []byte, password string) ([]byte, error) {
	key := c.DeriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

func (c *CryptoManager) Decrypt(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	key := c.DeriveKey(password)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}