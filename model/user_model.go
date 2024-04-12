package model

import (
	"time"
)

type UserModel struct {
	Id        string
	Fullname  string `json:"fullname" validate:"required,max=100"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	Role      RoleEnum
	Status    StatusEnum
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RegisterUserInput struct {
	Fullname string   `json:"fullname" validate:"required,max=100"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=6"`
	Role     RoleEnum `json:"role"`
}

type LoginUserInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdateUserInput struct {
	Fullname string     `json:"fullname" validate:"max=100"`
	Email    string     `json:"email" validate:"email"`
	Password string     `json:"password" `
	Role     RoleEnum   `json:"role"`
	Status   StatusEnum `json:"status"`
}

type UserFormatter struct {
	ID           string
	Fullname     string
	AccessToken  string
	RefreshToken string
	UserRole     RoleEnum
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}
