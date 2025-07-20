package main

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
)

const (
	notesFileName = ".notes.enc"
)

type Storage struct {
	filePath string
	crypto   *CryptoManager
}

func NewStorage() *Storage {
	homeDir, _ := os.UserHomeDir()
	filePath := filepath.Join(homeDir, notesFileName)
	
	return &Storage{
		filePath: filePath,
		crypto:   NewCryptoManager(),
	}
}

func (s *Storage) FileExists() bool {
	_, err := os.Stat(s.filePath)
	return err == nil
}

func (s *Storage) InitializeFile(password string) error {
	if s.FileExists() {
		return errors.New("notes file already exists")
	}

	if err := s.crypto.GenerateSalt(); err != nil {
		return err
	}

	emptyContent := []byte("")
	return s.SaveNotes(emptyContent, password)
}

func (s *Storage) SaveNotes(content []byte, password string) error {
	encrypted, err := s.crypto.Encrypt(content, password)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	
	if _, err := buf.Write(s.crypto.GetSalt()); err != nil {
		return err
	}
	
	if _, err := buf.Write(encrypted); err != nil {
		return err
	}

	return os.WriteFile(s.filePath, buf.Bytes(), 0600)
}

func (s *Storage) LoadNotes(password string) ([]byte, error) {
	if !s.FileExists() {
		return nil, errors.New("notes file does not exist")
	}

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return nil, err
	}

	if len(data) < saltSize {
		return nil, errors.New("invalid file format")
	}

	salt := data[:saltSize]
	encrypted := data[saltSize:]

	s.crypto.SetSalt(salt)

	decrypted, err := s.crypto.Decrypt(encrypted, password)
	if err != nil {
		return nil, errors.New("incorrect password or corrupted file")
	}

	return decrypted, nil
}

func (s *Storage) ChangePassword(oldPassword, newPassword string) error {
	content, err := s.LoadNotes(oldPassword)
	if err != nil {
		return err
	}

	if err := s.crypto.GenerateSalt(); err != nil {
		return err
	}

	return s.SaveNotes(content, newPassword)
}