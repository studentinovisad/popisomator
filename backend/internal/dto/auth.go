package dto

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_complexity"`
	FullName string `json:"full_name" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=admin manager user"`
	Status   string `validate:"required,oneof=requested active"`
}

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}
