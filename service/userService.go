package service

import (
	"context"
	"errors"
	"simple-product-api/models"
	"simple-product-api/utils"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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
		Role: user.Role,
	}
}

func ToUserResponse(user *models.User) (*models.UserResponse){
	return &models.UserResponse{
		Id: user.Id,
		Name: user.Name,
		Email: user.Email,
	}
}

func (us *UserService) Register(ctx context.Context, req *models.UserRequest) (*models.UserResponse, error) {

	//panggil fungsi hash password, hasilnya diset sebagai password data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	//Unique email check
	data, err := us.Repo.FindByEmail(ctx, req.Email)
	//either email gk ketemu (bagus) or mmg error waktu exec query
	if err != nil {
		return nil, err
	}
	if data != nil{
		return nil, errors.New("Email already exist!")
	}

	var product = &models.User{
		Id : uuid.New().String(),
		Name : req.Name,
		Password: string(hashedPassword),
		Email: req.Email,
		Role: utils.RoleUser,
	}

	data, err = us.Repo.Register(ctx, product)
	if err != nil {
		return nil, err
	}
	
	return ToUserResponse(data), nil
}

func (ur *UserService) Login(ctx context.Context, req *models.LoginRequest) (string, error){
	//Alur: get user from repo, bandingin password
	data, err := ur.Repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	//compare password
	err = bcrypt.CompareHashAndPassword([]byte(data.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}
	
	//placeholder untuk generate jwt token
	signedToken, err := utils.GenerateJWT(data.Id, data.Email, string(data.Role))
	if err != nil {
		return "", err
	}

	//return data dengan jwt
	return signedToken, err
}

func (us *UserService) GetAllUsers(ctx context.Context)([]*models.AdminUserResponse, error){
	var response []*models.AdminUserResponse
	
	data, err := us.Repo.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	for _, rows := range data{
		response = append(response, ToAdminUserResponse(rows))
	}
	return response, nil
}

func (us *UserService) GetUserById(ctx context.Context, id string) (*models.AdminUserResponse, error) {
	
	data, err := us.Repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToAdminUserResponse(data), nil
}

func (us *UserService) GetUserProfile(ctx context.Context, id string) (*models.UserResponse, error) {
	
	data, err := us.Repo.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return ToUserResponse(data), nil
}

func (us *UserService) UpdateUserProfile(ctx context.Context, id string, req *models.UserRequest) (*models.UserResponse, error) {
	//panggil fungsi hash password, hasilnya diset sebagai password data
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var data = &models.User{
		Name: req.Name,
		Password: string(hashedPassword),
		Email: req.Email,
	}

	data, err = us.Repo.UpdateUser(ctx, id, data)
	if err != nil {
		return nil, err
	}

	//if userRole = admin, toAdminUserResponse
	return ToUserResponse(data), nil
}

func (us *UserService) DeleteUser(ctx context.Context, id string) (*models.AdminUserResponse, error) {
	
	data, err := us.Repo.DeleteUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return ToAdminUserResponse(data), nil
}
