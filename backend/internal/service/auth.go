package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const sessionDuration = 24 * time.Hour

var hmacSecret []byte

func ConfigureJWT(secret string) {
	hmacSecret = []byte(secret)
}

// Returns JWT token if successful, error if not
func Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	if len(hmacSecret) == 0 {
		return "", errors.New("JWT secret is not configured")
	}

	user, err := db.Queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub": strconv.FormatInt(user.ID, 10),
		"exp": time.Now().Add(sessionDuration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

// Validates JWT token. If invalid, returns error. If valid, returns ID of user.
func ValidateToken(tokenStr string) (int64, error) {
	if len(hmacSecret) == 0 {
		return 0, errors.New("JWT secret is not configured")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return hmacSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return 0, err
	}

	subject, err := token.Claims.GetSubject()
	if err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func CreateUser(ctx context.Context, req dto.CreateUserRequest) (dto.User, error) {
	if err := dto.Validate(req); err != nil {
		return dto.User{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.User{}, err
	}

	user, err := db.Queries.CreateUser(ctx, repository.CreateUserParams{
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		Role:         repository.UserRole(req.Role),
	})
	if err != nil {
		return dto.User{}, err
	}

	return dto.ToUserDTO(user), nil
}
