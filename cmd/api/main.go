package main

import (
	"github.com/isaac-olle/gas-backoffice/cmd/api/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
