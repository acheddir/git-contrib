# git-contrib

A tool for analyzing Git commits and displaying a contribution graph on the command line.

## Features

- Display a contribution graph similar to GitHub's contribution calendar
- Filter contributions by email address
- Show commit counts or days of the month on the graph
- Use your own email from git config with the `--self` flag

## Installation

### Using Go

```bash
go install github.com/acheddir/git-contrib@latest
```

### Using Scoop (Windows)

You can install directly from the manifest:

```powershell
scoop install https://raw.githubusercontent.com/acheddir/git-contrib/main/git-contrib.json
```

## Usage

```bash
# Show contribution graph for all users in the current repository
git-contrib

# Show contribution graph for a specific email
git-contrib stats --email user@example.com

# Show contribution graph for your own commits (uses email from git config)
git-contrib stats --self

# Show commit counts on the graph
git-contrib stats --count

# Show days of the month on the graph
git-contrib stats --days
```

## Building from Source

### Linux/macOS

```bash
# Clone the repository
git clone https://github.com/acheddir/git-contrib.git
cd git-contrib

# Build
make
```

### Windows

```cmd
# Clone the repository
git clone https://github.com/acheddir/git-contrib.git
cd git-contrib

# Build using the provided batch file
build.bat

# Alternatively, you can build using Go directly
go mod tidy
go build -ldflags "-X github.com/acheddir/git-contrib/cmd.Version=0.1.0 -X github.com/acheddir/git-contrib/cmd.BuildDate=2023-05-13 -X github.com/acheddir/git-contrib/cmd.CommitHash=unknown"
```

## License

MIT
