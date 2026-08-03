package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/controller"
	"github.com/studentinovisad/popisomator/backend/internal/db"
)

func main() {
	// Initialise config
	if err := config.Init(); err != nil {
		log.Fatalf("unable to initialise config: %v", err)
	}

	// Connect to the database
	ctx := context.Background()
	pool, err := db.Connect(ctx, config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer pool.Close()

	// Routes
	http.HandleFunc("/ping", controller.Ping)
	http.HandleFunc("/health", controller.Healthcheck)
	http.HandleFunc("POST /auth/login", controller.Login)
	http.HandleFunc("POST /auth/logout", controller.Logout)
	http.HandleFunc("GET /user/details", controller.UserDetailsPersonal)

	// Listen for requests
	fmt.Println("Starting backend server on", config.CurrentConfig.Address)
	if err := http.ListenAndServe(config.CurrentConfig.Address, nil); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}
