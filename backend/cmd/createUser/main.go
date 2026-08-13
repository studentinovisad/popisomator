package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/studentinovisad/popisomator/backend/internal/config"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/service"
)

func main() {
	email := flag.String("email", "", "user email (required)")
	password := flag.String("password", "", "user password (required)")
	fullName := flag.String("full-name", "", "user full name (required)")
	role := flag.String("role", "user", "user role")
	flag.Parse()

	var missing []string
	if *email == "" {
		missing = append(missing, "email")
	}
	if *password == "" {
		missing = append(missing, "password")
	}
	if *fullName == "" {
		missing = append(missing, "full-name")
	}
	if len(missing) > 0 {
		fmt.Fprintf(os.Stderr, "missing required flag(s): %s\n\n", strings.Join(missing, ", "))
		flag.Usage()
		os.Exit(2)
	}

	if err := config.Init(); err != nil {
		log.Fatalf("unable to initialise config: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, config.CurrentConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer pool.Close()

	user, err := service.CreateUser(ctx, dto.CreateUserRequest{
		Email:    *email,
		Password: *password,
		FullName: *fullName,
		Role:     *role,
		Status:   "active",
	})
	if err != nil {
		log.Fatalf("unable to create user: %v", err)
	}

	fmt.Printf("created user %d (%s)\n", user.ID, user.Email)
}
