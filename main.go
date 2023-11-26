package main

import (
	"log"

	"github.com/whitfieldsdad/go-audit/cmd"
)

func main() {
	// Set process name to "go-audit"



	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
