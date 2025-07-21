# Note CLI - Secure Note-Taking Application

A simple command-line application for taking encrypted notes with customizable storage location, built with Go.

## Features

- Save encrypted notes with titles and timestamps
- Search notes by title
- List all note titles
- Password-protected encryption using AES-256-GCM
- Secure password handling with hidden input
- Customizable storage directory (e.g., Dropbox, Google Drive, etc.)
- Individual encrypted files for each note
- Cloud sync support through directory selection

## Installation

1. Clone the repository:
```bash
git clone https://github.com/user/note-cli.git
cd note-cli
```

2. Build the application:
```bash
go build -o note .
```

3. (Optional) Move to PATH:
```bash
sudo mv note /usr/local/bin/
```

## Usage

### First Time Setup
When you run the application for the first time, it will prompt you to choose where to store your notes:

```bash
$ note --list

First time setup - Please choose where to store your encrypted notes.
Examples: ~/Dropbox/, ~/Documents/notes/, ~/.my_notes/
Enter directory path: ~/Dropbox/my_notes/
Notes will be stored in: /home/user/Dropbox/my_notes/
Enter password: 
Available notes (0):
```

### Save a Note
Each note requires a title and content:
```bash
note --save --title "meeting_notes" "Team meeting at 3pm - discuss project timeline"
```

#### Adding Line Breaks to Notes
There are several ways to add line breaks to your notes:

**Method 1: Multiple Arguments (Recommended)**
Each argument becomes a separate line:
```bash
note --save --title "shopping_list" "Milk" "Bread" "Eggs" "Cheese"
```

**Method 2: Bash $'\n' Syntax**
Use bash's special quoting to convert \n to actual newlines:
```bash
note --save --title "todo" $'Task 1\nTask 2\nTask 3'
```

**Method 3: Command Substitution with echo -e**
```bash
note --save --title "formatted" "$(echo -e 'Line 1\nLine 2\nLine 3')"
```

**Method 4: Reading from File**
```bash
note --save --title "github_guide" $(cat Github.txt)
```

**Note:** Literal `\n` characters will appear as text, not line breaks.

### Open and Read a Note
Open and display a note by its title:
```bash
note --open --title "meeting_notes"
```

### Find Notes by Title
Search for notes by title (case-insensitive):
```bash
note --find "meeting"
```

### List All Notes
Display all available note titles:
```bash
note --list
```

### Help
```bash
note --help
```

## Storage Configuration

After first setup, your chosen directory is saved in `~/.note-cli-config.json`:

```json
{
  "notes_directory": "/home/user/Dropbox/my_notes/"
}
```

### Cloud Sync Examples
- **Dropbox**: `~/Dropbox/notes/`
- **Google Drive**: `~/Google Drive/notes/`
- **OneDrive**: `~/OneDrive/notes/`
- **Local**: `~/Documents/notes/` or `~/.my_notes/`

## Security

- Notes are encrypted using AES-256-GCM
- Password-based key derivation using PBKDF2 with 100,000 iterations
- Each note is stored as an individual encrypted `.enc` file
- File permissions set to 0600 (read/write for owner only)
- Passwords are never stored, only used for key derivation
- Each note can use a different password

## File Structure

```
chosen-directory/
├── note_title_1.enc
├── shopping_list.enc
├── meeting_notes.enc
└── important_todo.enc
```

Each `.enc` file contains:
- Salt (32 bytes)
- Encrypted content (timestamp + note content)

## Technical Details

- Written in Go
- Uses `golang.org/x/term` for secure password input
- Uses `golang.org/x/crypto/pbkdf2` for key derivation
- Each note is timestamped with format: `[YYYY-MM-DD HH:MM:SS]`
- Title sanitization (spaces and slashes converted to underscores)

## Example Workflow

```bash
# First time - choose directory
$ note --save --title "first_note" "Hello world"
First time setup - Please choose where to store your encrypted notes.
Examples: ~/Dropbox/, ~/Documents/notes/, ~/.my_notes/
Enter directory path: ~/Dropbox/notes/
Notes will be stored in: /home/user/Dropbox/notes/
Enter password: 
Note saved successfully with title: first_note

# Subsequent uses - no setup needed
$ note --save --title "shopping" "Buy milk and eggs"
Enter password: 
Note saved successfully with title: shopping

$ note --list
Enter password: 
Available notes (2):
  - first_note
  - shopping

$ note --find "shop"
Enter password: 
Found 1 note(s) with matching title:
  - shopping

$ note --open --title "shopping"
Enter password: 
Title: shopping
Content:
[2023-12-20 15:30:45] Buy milk and eggs
```

## License

MIT License