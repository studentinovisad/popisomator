package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/controller"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

func main() {
	// Initialise config
	if err := config.Init(); err != nil {
		log.Fatalf("unable to initialise config: %v", err)
	}

	// Connect to the database
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("unable to create postgres connection pool: %v", err)
	}
	defer pool.Close()

	queries := repository.New(pool)
	healthController := controller.NewHealthController(queries)

	// Routes
	http.HandleFunc("/ping", controller.Ping)
	http.HandleFunc("/health", healthController.Healthcheck)

	// Listen for requests
	fmt.Println("Starting backend server on", config.CurrentConfig.Address)
	if err := http.ListenAndServe(config.CurrentConfig.Address, nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
