package service

import (
	"context"

	"github.com/studentinovisad/popisomator/backend/internal/db"
	"github.com/studentinovisad/popisomator/backend/internal/dto"
	"github.com/studentinovisad/popisomator/backend/internal/repository"
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

func ListUsers(ctx context.Context, limit, offset int32) (dto.UsersPage, error) {
	users, err := db.Queries.ListUsers(ctx, repository.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return dto.UsersPage{}, err
	}

	total, err := db.Queries.CountUsers(ctx)
	if err != nil {
		return dto.UsersPage{}, err
	}

	items := make([]dto.User, len(users))
	for index, user := range users {
		items[index] = dto.ToUserDTO(user)
	}

	return dto.UsersPage{
		Items:  items,
		Limit:  limit,
		Offset: offset,
		Total:  total,
	}, nil
}

func UpdateUserRole(ctx context.Context, id int64, req dto.UpdateRoleRequest) (dto.User, error) {
	if err := dto.Validate(req); err != nil {
		return dto.User{}, err
	}

	user, err := db.Queries.UpdateRole(ctx, repository.UpdateRoleParams{
		ID:   id,
		Role: repository.UserRole(req.Role),
	})
	if err != nil {
		return dto.User{}, err
	}

	return dto.ToUserDTO(user), nil
}
