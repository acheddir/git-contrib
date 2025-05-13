package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Print(time.Now().UTC().Format("2006-01-02"))
}
