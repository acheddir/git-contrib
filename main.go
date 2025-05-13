package main

import (
	"github.com/acheddir/git-contrib/cmd"
)

func main() {
	cmd.Execute()
}

type Person struct {
	Name string
}

func (p *Person) PrintScore() string {
	return "Score: 100"
}

func (p *Person) delete() bool {
	return true
}

var isDeletable interface {
	Delete() bool
}
