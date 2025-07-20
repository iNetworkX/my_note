package main

import (
	"fmt"
	"strings"
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

func (n *NotesManager) Initialize(password string) error {
	return n.storage.InitializeFile(password)
}

func (n *NotesManager) IsInitialized() bool {
	return n.storage.FileExists()
}

func (n *NotesManager) AddNote(note string, password string) error {
	content, err := n.storage.LoadNotes(password)
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	newNote := fmt.Sprintf("[%s] %s\n", timestamp, note)

	updatedContent := append(content, []byte(newNote)...)
	return n.storage.SaveNotes(updatedContent, password)
}

func (n *NotesManager) FindNotes(searchTerm string, password string) ([]string, error) {
	content, err := n.storage.LoadNotes(password)
	if err != nil {
		return nil, err
	}

	allNotes := strings.Split(string(content), "\n")
	var matchingNotes []string

	searchLower := strings.ToLower(searchTerm)
	for _, note := range allNotes {
		if note != "" && strings.Contains(strings.ToLower(note), searchLower) {
			matchingNotes = append(matchingNotes, note)
		}
	}

	return matchingNotes, nil
}

func (n *NotesManager) ListAllNotes(password string) ([]string, error) {
	content, err := n.storage.LoadNotes(password)
	if err != nil {
		return nil, err
	}

	allNotes := strings.Split(string(content), "\n")
	var notes []string

	for _, note := range allNotes {
		if note != "" {
			notes = append(notes, note)
		}
	}

	return notes, nil
}

func (n *NotesManager) ChangePassword(oldPassword, newPassword string) error {
	return n.storage.ChangePassword(oldPassword, newPassword)
}