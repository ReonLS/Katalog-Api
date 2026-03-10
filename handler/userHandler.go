package handler

import (
	"encoding/json"
	"net/http"
	"simple-product-api/models"
	"simple-product-api/service"
	"simple-product-api/utils"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

// @Summary Register a new user account
// @description Create a new user account with the User role
// @tags Public
// @accept json
// @Produce json
// @Param user body models.UserRequest true "User credentials"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /register [POST]
func (uh *UserHandler) Register(rw http.ResponseWriter, r *http.Request) {
	req := &models.UserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := utils.ValidateRequest(req.Name, req.Email, req.Password); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}
		GenerateError(rw, strings.Join(joinedError, " ,"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.Register(r.Context(), req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusCreated, response)
}

// @Summary Authenticate user
// @description Authenticate user credentials and return a signed JWT token
// @tags Public
// @accept json
// @Produce plain
// @Param user body models.LoginRequest true "User credentials"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} models.BadRequestResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /login [POST]
func (uh *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	var req = &models.LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := utils.ValidateLogin(req.Email, req.Password); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}
		GenerateError(rw, strings.Join(joinedError, " ,"), http.StatusBadRequest)
		return
	}

	token, err := uh.Service.Login(r.Context(), req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, token)
}

// @Summary Get user profile
// @description Returns user's credential information
// @tags User
// @accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /user [GET]
// @Security BearerAuth
func (uh *UserHandler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Failed Claims", http.StatusUnauthorized)
		return
	}

	response, err := uh.Service.GetUserProfile(r.Context(), claims.Id)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Update user profile
// @description Returns user's newly updated credentials
// @tags User
// @accept json
// @Produce json
// @Param user body models.UserRequest true "User credentials"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /user [PUT]
// @Security BearerAuth
func (uh *UserHandler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	var req = &models.UserRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Failed Claims", http.StatusUnauthorized)
		return
	}

	if err := utils.ValidateRequest(req.Name, req.Email, req.Password); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}
		GenerateError(rw, strings.Join(joinedError, " ,"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.UpdateUserProfile(r.Context(), claims.Id, req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Get all user
// @description Return all user's profile information
// @tags Admin
// @accept json
// @Produce json
// @Success 200 {array} models.AdminUserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/user [GET]
// @Security BearerAuth
func (uh *UserHandler) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	response, err := uh.Service.GetAllUsers(r.Context())
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Get user profile
// @description Returns a user's profile by userID
// @tags Admin
// @accept json
// @Produce json
// @Param id path string true "UserID"
// @Success 200 {object} models.AdminUserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /admin/user/{id} [GET]
// @Security BearerAuth
func (uh *UserHandler) AdminGetUserProfile(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return 
	}

	response, err := uh.Service.GetUserById(r.Context(), userID)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Delete user
// @description Returns deleted user's profile information
// @tags Admin
// @accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.AdminUserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/user/{id} [DELETE]
// @Security BearerAuth
func (uh *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return 
	}

	response, err := uh.Service.DeleteUser(r.Context(), userID)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}
