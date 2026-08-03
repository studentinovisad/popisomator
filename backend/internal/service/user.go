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
