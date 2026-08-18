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

func ListUsers(ctx context.Context, request dto.ListUsersRequest) (dto.UsersPage, error) {
	users, err := db.Queries.ListUsers(ctx, repository.ListUsersParams{
		Search:       request.Search,
		RoleFilter:   request.Role,
		StatusFilter: request.Status,
		PageOffset:   request.Offset,
		PageLimit:    request.Limit,
	})
	if err != nil {
		return dto.UsersPage{}, err
	}

	total, err := db.Queries.CountUsers(ctx, repository.CountUsersParams{
		Search:       request.Search,
		RoleFilter:   request.Role,
		StatusFilter: request.Status,
	})
	if err != nil {
		return dto.UsersPage{}, err
	}

	items := make([]dto.User, len(users))
	for index, user := range users {
		items[index] = dto.ToUserDTO(user)
	}

	return dto.UsersPage{
		Items:  items,
		Limit:  request.Limit,
		Offset: request.Offset,
		Total:  total,
	}, nil
}

func UpdateUser(ctx context.Context, id int64, req dto.UpdateUserRequest) (dto.User, error) {
	if err := dto.Validate(req); err != nil {
		return dto.User{}, err
	}

	var user repository.User
	if req.Role != nil || req.Status != nil {
		tx, err := db.BeginTransaction(ctx)
		if err != nil {
			return dto.User{}, err
		}
		defer tx.Rollback(ctx)
		queriesTx := db.Queries.WithTx(tx)

		if req.Role != nil {
			var err error
			if user, err = queriesTx.UpdateUserRole(ctx, repository.UpdateUserRoleParams{
				ID:   id,
				Role: repository.UserRole(*req.Role),
			}); err != nil {
				return dto.User{}, err
			}
		}

		if req.Status != nil {
			var err error
			if user, err = queriesTx.UpdateUserStatus(ctx, repository.UpdateUserStatusParams{
				ID:     id,
				Status: repository.UserStatus(*req.Status),
			}); err != nil {
				return dto.User{}, err
			}
		}

		if err := tx.Commit(ctx); err != nil {
			return dto.User{}, err
		}
	} else {
		return dto.User{}, ErrNoUpdateFields
	}

	return dto.ToUserDTO(user), nil
}

func DeleteUser(ctx context.Context, id int64) error {
	rowsAffected, err := db.Queries.DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
