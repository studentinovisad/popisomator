package dto

import "github.com/studentinovisad/popisomator/backend/internal/repository"

type User struct {
	ID       int64               `json:"id"`
	Email    string              `json:"email"`
	FullName string              `json:"full_name"`
	Role     repository.UserRole `json:"role"`
}

type UpdateRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin manager user"`
}

type UsersPage struct {
	Items  []User `json:"items"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	Total  int64  `json:"total"`
}

func ToUserDTO(user repository.User) User {
	return User{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}
}
