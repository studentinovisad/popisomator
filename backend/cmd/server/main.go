package main

import (
	"fmt"
	"net/http"

	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/controller"
)

func main() {
	// Initialise config
	config.Init()

	// Routes
	http.HandleFunc("/ping", controller.Ping)

	// Listen for requests
	fmt.Println("Starting backend server on", config.CurrentConfig.Address)
	http.ListenAndServe(config.CurrentConfig.Address, nil)
}
