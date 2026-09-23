package dto

import (
	"github.com/studentinovisad/popisomator/backend/internal/repository"
)

type User struct {
	ID       int64                 `json:"id"`
	Email    string                `json:"email"`
	FullName string                `json:"full_name"`
	Role     repository.UserRole   `json:"role"`
	Status   repository.UserStatus `json:"status"`
}

type UpdateUserRequest struct {
	Email    *string `json:"email" validate:"omitempty,email"`
	FullName *string `json:"full_name" validate:"omitempty,max=200"`
	Role     *string `json:"role" validate:"omitempty,oneof=admin manager user"`
	Status   *string `json:"status" validate:"omitempty,oneof=active requested"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8,password_complexity"`
}

type UsersPage struct {
	Items  []User `json:"items"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Total  int64  `json:"total"`
}

type ListUsersRequest struct {
	Limit  int32
	Offset int32
	Search string
	Role   *repository.UserRole
	Status *repository.UserStatus
}

func ToUserDTO(user repository.User) User {
	return User{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Status:   user.Status,
	}
}
