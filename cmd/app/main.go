package main

import (
	"log"
	"os"

	"github.com/joalvarez/toolkit/internal/ui"
)

func main() {
	if err := ui.Run(); err != nil {
		log.New(os.Stderr, "", 0).Fatal(err)
	}
}
