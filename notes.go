package main

import (
	"fmt"
	"time"
)

type NotesManager struct {
	storage *Storage
}

func NewNotesManager() *NotesManager {
	return &NotesManager{
		storage: NewStorage(),
	}
}

func (n *NotesManager) IsInitialized() bool {
	return n.storage.IsInitialized()
}

func (n *NotesManager) AddNote(title string, note string, password string) error {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	content := fmt.Sprintf("[%s]\n%s", timestamp, note)
	
	return n.storage.SaveNote(title, []byte(content), password)
}

func (n *NotesManager) FindNotes(titleSearch string, password string) ([]string, error) {
	titles, err := n.storage.ListTitles()
	if err != nil {
		return nil, err
	}

	var matchingTitles []string
	for _, title := range titles {
		// Case-insensitive title search
		if containsIgnoreCase(title, titleSearch) {
			matchingTitles = append(matchingTitles, title)
		}
	}

	return matchingTitles, nil
}

func (n *NotesManager) ListAllNotes(password string) ([]string, error) {
	return n.storage.ListTitles()
}

func (n *NotesManager) GetNoteContent(title string, password string) (string, error) {
	content, err := n.storage.LoadNote(title, password)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (n *NotesManager) DeleteNote(title string, password string) error {
	// Verify password by trying to load the note first
	_, err := n.storage.LoadNote(title, password)
	if err != nil {
		return err
	}
	return n.storage.DeleteNote(title)
}

// Helper function for case-insensitive string contains
func containsIgnoreCase(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return contains(s, substr)
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if 'A' <= c && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}