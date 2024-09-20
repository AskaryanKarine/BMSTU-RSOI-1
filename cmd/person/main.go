package main

import (
	"github.com/AskaryanKarine/BMSTU-ds-1/internal/app"
	"github.com/charmbracelet/log"
)

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	if err = a.Start(); err != nil {
		log.Fatal(err)
	}
}
