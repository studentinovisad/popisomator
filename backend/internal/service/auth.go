package service

import (
	"context"
	"crypto/rand"

	"github.com/golang-jwt/jwt/v5"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"golang.org/x/crypto/bcrypt"
)

var hmacSecret string = rand.Text()

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Returns JWT token if successful, error if not
func Login(ctx context.Context, req LoginRequest) (string, error) {
	user, err := db.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":   user.ID,
		"role": user.Role,
	})

	tokenStr, err := token.SignedString([]byte(hmacSecret))
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
