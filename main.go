package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
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
		save     = flag.Bool("save", false, "Save a new note")
		open     = flag.Bool("open", false, "Open and read a note")
		title    = flag.String("title", "", "Title for the note (required with --save or --open)")
		find     = flag.String("find", "", "Find notes by title")
		list     = flag.Bool("list", false, "List all note titles")
		password = flag.Bool("password", false, "Show password info (deprecated)")
		help     = flag.Bool("help", false, "Show help")
	)

	flag.Parse()

	if *help || (!*list && !*save && *find == "" && !*password && !*open) {
		fmt.Println("Note CLI - Secure note-taking application")
		fmt.Println("\nUsage:")
		fmt.Println("  note -save -title \"my_note\" \"Your note content\"     Save a new note")
		fmt.Println("  note -open -title \"my_note\"                         Open and read a note")
		fmt.Println("  note -find \"title\"                                   Find notes by title")
		fmt.Println("  note -list                                           List all note titles")
		fmt.Println("  note -password                                       Show password info (deprecated)")
		fmt.Println("  note -help                                           Show this help")
		fmt.Println("  note -h                                              Show flag usage")
		os.Exit(0)
	}

	notes := NewNotesManager()

	if *password {
		fmt.Println("Password management is not needed with the new title-based system.")
		fmt.Println("Each note can use a different password.")
		return
	}

	pass, err := readPassword("Enter password: ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading password: %v\n", err)
		os.Exit(1)
	}

	switch {
	case *save:
		if *title == "" {
			fmt.Fprintf(os.Stderr, "Error: -title is required when using -save\n")
			os.Exit(1)
		}
		
		// Get content from remaining arguments
		args := flag.Args()
		if len(args) == 0 {
			fmt.Fprintf(os.Stderr, "Error: content is required when using -save\n")
			fmt.Fprintf(os.Stderr, "Usage: note -save -title \"your_title\" \"your content here\"\n")
			os.Exit(1)
		}
		content := strings.Join(args, "\n")
		
		if err := notes.AddNote(*title, content, pass); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving note: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Note saved successfully with title: %s\n", *title)

	case *open:
		if *title == "" {
			fmt.Fprintf(os.Stderr, "Error: -title is required when using -open\n")
			os.Exit(1)
		}
		content, err := notes.GetNoteContent(*title, pass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening note: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Title: %s\n", *title)
		fmt.Printf("Content:\n%s\n", content)

	case *find != "":
		results, err := notes.FindNotes(*find, pass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error searching notes: %v\n", err)
			os.Exit(1)
		}
		if len(results) == 0 {
			fmt.Printf("No notes found with title containing: %s\n", *find)
		} else {
			fmt.Printf("Found %d note(s) with matching title:\n", len(results))
			for _, title := range results {
				fmt.Printf("  - %s\n", title)
			}
		}

	case *list:
		titles, err := notes.ListAllNotes(pass)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error listing notes: %v\n", err)
			os.Exit(1)
		}
		if len(titles) == 0 {
			fmt.Println("No notes found.")
		} else {
			fmt.Printf("Available notes (%d):\n", len(titles))
			for _, title := range titles {
				fmt.Printf("  - %s\n", title)
			}
		}
	}
}