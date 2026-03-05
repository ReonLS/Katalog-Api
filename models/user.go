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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AdminUserResponse struct {
	Id    string     `json:"id"`
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Role  utils.Role `json:"role"`
}
