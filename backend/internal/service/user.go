package service

import (
	"context"

	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
)

func GetUserDetails(ctx context.Context, id int64) (dto.User, error) {
	user, err := db.Queries.GetUserByID(ctx, id)
	if err != nil {
		return dto.User{}, err
	}

	userDTO := dto.ToUserDTO(user)

	return userDTO, nil
}

func GetUserByEmail(ctx context.Context, email string) (dto.User, error) {
	user, err := db.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		return dto.User{}, err
	}

	userDTO := dto.ToUserDTO(user)
	return userDTO, nil
}
