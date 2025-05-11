package main

import (
	"flag"
	"git-contrib/pkg/commands"
)

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&email, "email", "your@email.com", "the email to scan")
	flag.Parse()

	if folder != "" {
		commands.Scan(folder)
		return
	}

	commands.Stats(email)
}
