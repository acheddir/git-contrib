package main

import (
	"flag"
	"git-contrib/pkg/commands"
)

func main() {
	var folder string
	var email string

	flag.StringVar(&folder, "add", "", "add a new folder to scan for Git repositories")
	flag.StringVar(&email, "email", "acheddir@redsen.ch", "the email to scan")
	flag.Parse()

	if folder != "" {
		commands.Scan(folder)
		return
	}

	err := commands.Stats(email)
	if err != nil {
		return
	}
}
