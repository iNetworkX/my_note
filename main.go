package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

	"golang.org/x/term"
)

func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

func main() {
	var (
		save     = flag.String("save", "", "Save a new note")
		find     = flag.String("find", "", "Find notes containing text")
		list     = flag.Bool("list", false, "List all notes")
		password = flag.Bool("password", false, "Set or change password")
		help     = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help || (!*list && *save == "" && *find == "" && !*password) {
		fmt.Println("Note CLI - Secure note-taking application")
		fmt.Println("\nUsage:")
		fmt.Println("  note --save \"Your note content\"    Save a new note")
		fmt.Println("  note --find \"search term\"          Find notes containing text")
		fmt.Println("  note --list                        List all notes")
		fmt.Println("  note --password                    Set or change password")
		fmt.Println("  note --help                        Show this help")
		os.Exit(0)
	}

	notes := NewNotesManager()

	if *password {
		handlePasswordChange(notes)
		return
	}

	if !notes.IsInitialized() {
		fmt.Println("First time setup - please set a password for your notes.")
		if err := setupInitialPassword(notes); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Password set successfully. You can now start saving notes.")
		return
	}

	pass, err := readPassword("Enter password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *save != "":
		if err := notes.AddNote(*save, pass); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Note saved successfully.")

	case *find != "":
		results, err := notes.FindNotes(*find, pass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching notes: %v\n", err)
			os.Exit(1)
		}
		if len(results) == 0 {
			fmt.Println("No notes found containing:", *find)
		} else {
			fmt.Printf("Found %d note(s):\n", len(results))
			for _, note := range results {
				fmt.Println(note)
			}
		}

	case *list:
		allNotes, err := notes.ListAllNotes(pass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing notes: %v\n", err)
			os.Exit(1)
		}
		if len(allNotes) == 0 {
			fmt.Println("No notes found.")
		} else {
			fmt.Printf("Total notes: %d\n", len(allNotes))
			for _, note := range allNotes {
				fmt.Println(note)
			}
		}
	}
}

func setupInitialPassword(notes *NotesManager) error {
	pass1, err := readPassword("Enter new password: ")
	if err != nil {
		return err
	}

	pass2, err := readPassword("Confirm password: ")
	if err != nil {
		return err
	}

	if pass1 != pass2 {
		return fmt.Errorf("passwords do not match")
	}

	if len(pass1) < 4 {
		return fmt.Errorf("password must be at least 4 characters long")
	}

	return notes.Initialize(pass1)
}

func handlePasswordChange(notes *NotesManager) {
	if !notes.IsInitialized() {
		if err := setupInitialPassword(notes); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Password set successfully.")
		return
	}

	oldPass, err := readPassword("Enter current password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	newPass1, err := readPassword("Enter new password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	newPass2, err := readPassword("Confirm new password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	if newPass1 != newPass2 {
		fmt.Fprintf(os.Stderr, "Error: passwords do not match\n")
		os.Exit(1)
	}

	if len(newPass1) < 4 {
		fmt.Fprintf(os.Stderr, "Error: password must be at least 4 characters long\n")
		os.Exit(1)
	}

	if err := notes.ChangePassword(oldPass, newPass1); err != nil {
		fmt.Fprintf(os.Stderr, "Error changing password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password changed successfully.")
}