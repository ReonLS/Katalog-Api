package service

import (
	"simple-product-api/models"
)

type UserService struct{
	Repo models.UserRepository
}

//constructors
func NewUserService(repo models.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func ToAdminUserResponse(user *models.User) (*models.AdminUserResponse){
	return &models.AdminUserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
		IsAdmin: user.IsAdmin,
	}
}

func ToUserResponse(user *models.User) (*models.UserResponse){
	return &models.UserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
	}
}

func (us *UserService) GetAllUsers()([]*models.AdminUserResponse, error){
	var response []*models.AdminUserResponse
	
	data, err := us.Repo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	for _, rows := range data{
		response = append(response, ToAdminUserResponse(rows))
	}
	return response, nil
}

func (us *UserService) AdminGetUserbyId(id int) (*models.AdminUserResponse, error) {
	
	data, err := us.Repo.GetUserbyId(id)
	if err != nil {
		return nil, err
	}

	return ToAdminUserResponse(data), nil
}

func (us *UserService) GetUserbyId(id int) (*models.UserResponse, error) {
	
	data, err := us.Repo.GetUserbyId(id)
	if err != nil {
		return nil, err
	}

	return ToUserResponse(data), nil
}

func (us *UserService) CreateUser(req *models.UserRequest) (*models.AdminUserResponse, error) {
	
	var data = &models.User{
		Name: req.Name,
		Email: req.Email,
		IsAdmin: false,
	}

	data, err := us.Repo.CreateUser(data)
	if err != nil {
		return nil, err
	}

	return ToAdminUserResponse(data), nil
}

func (us *UserService) UpdateUser(id int, req *models.UserRequest) (*models.UserResponse, error) {
	
	var data = &models.User{
		Name: req.Name,
		Email: req.Email,
	}

	data, err := us.Repo.UpdateUser(id, data)
	if err != nil {
		return nil, err
	}

	return ToUserResponse(data), nil
}

func (us *UserService) DeleteUser(id int) (*models.AdminUserResponse, error) {
	
	data, err := us.Repo.DeleteUser(id)
	if err != nil {
		return nil, err
	}

	return ToAdminUserResponse(data), nil
}
