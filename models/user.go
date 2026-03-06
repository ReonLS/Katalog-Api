package models

import (
	"context"
	"simple-product-api/utils"
)

// Repository Pattern
type UserRepository interface {
	Register(ctx context.Context, model *User) (*User, error)
	FindByEmail(ttx context.Context, email string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)                      //admin
	GetUserById(ctx context.Context, id string) (*User, error)             //admin
	UpdateUser(ctx context.Context, id string, model *User) (*User, error) //both admin and user
	DeleteUser(ctx context.Context, id string) (*User, error)              //admin
}

type User struct {
	Id       string
	Name     string
	Password string
	Email    string
	Role     utils.Role
}

// Update & Register
type UserRequest struct {
	Name     string `json:"name" example:"John"`
	Email    string `json:"email" example:"john@gmail.com"`
	Password string `json:"password" example:"john1234"`
}

type LoginRequest struct {
	Email    string `json:"email" example:"john@gmail.com"`
	Password string `json:"password" example:"john1234"`
}

type UserResponse struct {
	Id    string `json:"id" example:"4dc0e876-8f8a-41c4-be90-e35da0105ccf"`
	Name     string `json:"name" example:"John"`
	Email    string `json:"email" example:"john@gmail.com"`
}

type AdminUserResponse struct {
	Id    string `json:"id" example:"4dc0e876-8f8a-41c4-be90-e35da0105ccf"`
	Name     string `json:"name" example:"John"`
	Email    string `json:"email" example:"john@gmail.com"`
	Role  utils.Role `json:"role" example:"User"`
}
