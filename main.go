package main

import (
	"log"

	"github.com/whitfieldsdad/go-audit/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
