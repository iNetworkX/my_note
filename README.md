# Note CLI - Secure Note-Taking Application

A simple command-line application for taking encrypted notes, built with Go.

## Features

- Save encrypted notes with timestamps
- Search notes by content
- List all notes
- Password-protected encryption using AES-256-GCM
- Secure password handling with hidden input
- Single encrypted file storage

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
When you run the application for the first time, it will prompt you to set a password:
```bash
note --password
```

### Save a Note
```bash
note --save "Meeting at 3pm with the team"
```

### Find Notes
Search for notes containing specific text:
```bash
note --find "meeting"
```

### List All Notes
```bash
note --list
```

### Change Password
```bash
note --password
```

### Help
```bash
note --help
```

## Security

- Notes are encrypted using AES-256-GCM
- Password-based key derivation using PBKDF2 with 100,000 iterations
- Encrypted notes are stored in `~/.notes.enc`
- File permissions set to 0600 (read/write for owner only)
- Passwords are never stored, only used for key derivation

## Technical Details

- Written in Go
- Uses `golang.org/x/term` for secure password input
- Uses `golang.org/x/crypto/pbkdf2` for key derivation
- Each note is timestamped with format: `[YYYY-MM-DD HH:MM:SS] note content`

## License

MIT License# my_note
