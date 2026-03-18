package main

import (
	"log"

	"github.com/romangurevitch/go-training/internal/bank/app"
	"github.com/romangurevitch/go-training/internal/bank/config"
)

func init() {
	config.Init()
}

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("application failed: %v", err)
	}
}
