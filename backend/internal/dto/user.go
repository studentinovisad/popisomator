package dto

import "github.com/studentinovisad/popisomator/backend/internal/repository"

type User struct {
	ID       int64                 `json:"id"`
	Email    string                `json:"email"`
	FullName string                `json:"full_name"`
	Role     repository.UserRole   `json:"role"`
	Status   repository.UserStatus `json:"status"`
}

type UpdateUserRequest struct {
	Role   *string `json:"role" validate:"omitempty,oneof=admin manager user"`
	Status *string `json:"status" validate:"omitempty,oneof=active requested"`
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
	Role   string
	Status string
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
