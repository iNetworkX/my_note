package main

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	dirPath string
	crypto  *CryptoManager
	config  *ConfigManager
}

func NewStorage() *Storage {
	config := NewConfigManager()
	
	// Get configured directory or set it up
	var dirPath string
	if config.Exists() {
		dirPath, _ = config.GetNotesDirectory()
	}
	
	return &Storage{
		dirPath: dirPath,
		crypto:  NewCryptoManager(),
		config:  config,
	}
}

func (s *Storage) IsInitialized() bool {
	return s.config.Exists()
}

func (s *Storage) EnsureInitialized() error {
	if s.dirPath == "" {
		dirPath, err := s.config.SetupNotesDirectory()
		if err != nil {
			return err
		}
		s.dirPath = dirPath
	}
	return nil
}

func (s *Storage) ensureDir() error {
	if err := s.EnsureInitialized(); err != nil {
		return err
	}
	return os.MkdirAll(s.dirPath, 0700)
}

func (s *Storage) getFilePath(title string) string {
	// Sanitize title to be a valid filename
	safeTitle := strings.ReplaceAll(title, "/", "_")
	safeTitle = strings.ReplaceAll(safeTitle, " ", "_")
	return filepath.Join(s.dirPath, safeTitle+".enc")
}

func (s *Storage) FileExists(title string) bool {
	_, err := os.Stat(s.getFilePath(title))
	return err == nil
}

func (s *Storage) SaveNote(title string, content []byte, password string) error {
	if err := s.ensureDir(); err != nil {
		return fmt.Errorf("failed to create notes directory: %v", err)
	}

	// Generate new salt for each note
	if err := s.crypto.GenerateSalt(); err != nil {
		return err
	}

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

	return os.WriteFile(s.getFilePath(title), buf.Bytes(), 0600)
}

func (s *Storage) LoadNote(title string, password string) ([]byte, error) {
	if !s.FileExists(title) {
		return nil, fmt.Errorf("note '%s' does not exist", title)
	}

	data, err := os.ReadFile(s.getFilePath(title))
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

func (s *Storage) ListTitles() ([]string, error) {
	if err := s.ensureDir(); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(s.dirPath)
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".enc") {
			title := strings.TrimSuffix(entry.Name(), ".enc")
			titles = append(titles, title)
		}
	}

	return titles, nil
}

func (s *Storage) DeleteNote(title string) error {
	if !s.FileExists(title) {
		return fmt.Errorf("note '%s' does not exist", title)
	}
	return os.Remove(s.getFilePath(title))
}