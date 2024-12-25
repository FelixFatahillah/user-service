package dtos

import "user-service/internal/domain/user/models"

type CreateUserDto struct {
	FirstName   string          `json:"first_name" validate:"required"`
	LastName    *string         `json:"last_name" validate:"required"`
	PhoneNumber *string         `json:"phone_number" validate:"required"`
	Email       string          `json:"email" validate:"required"`
	Password    string          `json:"password" validate:"required"`
	Role        models.RoleEnum `json:"role" validate:"required"`
	IsActive    bool            `json:"is_active"`
}

type CreateUserResponseDto struct {
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       string  `json:"email"`
}

type UpdateUserDto struct {
	ID          string  `json:"id"`
	FirstName   string  `json:"first_name"`
	LastName    *string `json:"last_name"`
	PhoneNumber *string `json:"phone_number"`
	Email       string  `json:"email"`
	Password    string  `json:"password"`
	IsActive    bool    `json:"is_active"`
}
