package models

//Repository Pattern
type UserRepository interface{
	GetAllUsers()([]*User, error)
	GetUserbyId(id int)(*User, error)
	CreateUser(model *UserRequest)(*User, error)
	UpdateUser(id int, req *UserRequest)(*User, error)
	DeleteUser(id int)(*User, error)
}

type User struct{
	Id int
	Name string
	Email string
	IsAdmin bool
}

//hanya admin yg bisa request user (PUSH)
type UserRequest struct{
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type UserResponse struct{
	Id int `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

type AdminUserResponse struct{
	Id int `json:"id" binding:"required"`
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
	IsAdmin bool `json:"isAdmin" binding:"required"`
}