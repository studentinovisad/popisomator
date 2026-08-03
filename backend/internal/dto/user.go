package dto

import "github.com/studentinovisad/popisomator/backend/internal/repository"

type User struct {
	ID       int64               `json:"id"`
	Email    string              `json:"email"`
	FullName string              `json:"full_name"`
	Role     repository.UserRole `json:"role"`
}

func ToUserDTO(user repository.User) User {
	return User{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
	}
}
