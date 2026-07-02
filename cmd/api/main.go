package main

import (
	"github.com/isaac-olle/gas/cmd/api/bootstrap"
	"log"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
