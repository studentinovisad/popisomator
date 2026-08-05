package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/router"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func main() {
	// Initialise config
	if err := config.Init(); err != nil {
		log.Fatalf("unable to initialise config: %v", err)
	}
	if config.CurrentConfig.JWTSecret == "" {
		log.Fatal("POPISOMATOR_JWT_SECRET is required")
	}
	service.ConfigureJWT(config.CurrentConfig.JWTSecret)

	// Connect to the database
	ctx := context.Background()
	pool, err := db.Connect(ctx, config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer pool.Close()

	mux := router.New()

	// Listen for requests
	address := config.CurrentConfig.Address()
	fmt.Println("Starting backend server on", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
