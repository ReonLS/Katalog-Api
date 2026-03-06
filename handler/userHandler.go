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

// @Summary Register account
// @description Generate a user account when successful
// @tags Public
// @accept json
// @Produce json
// @Param user body models.UserRequest true "Create Account"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /register [POST]
func (uh *UserHandler) Register(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var req = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//validate
	if err := utils.ValidateRequest(req.Name, req.Email, req.Password); len(err) > 0 {
		//Access setiap error, join ke joinedError, return sebagai message
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		GenerateError(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.Register(r.Context(), req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Log in
// @description Authenticate user to generate JWT
// @tags Public
// @accept json
// @Produce plain
// @Param user body models.LoginRequest true "Login Account"
// @Success 200 {string} string "JWT Token"
// @Failure 400 {object} models.BadRequestResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /login [POST]
func (uh *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var req = &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		GenerateError(rw, "Error Request", http.StatusBadRequest)
		return
	}

	//validate
	if err := utils.ValidateLogin(req.Email, req.Password); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		GenerateError(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	//could be either error or token
	token, err := uh.Service.Login(r.Context(), req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(token)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Get profile
// @description User get their profile
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
	rw.Header().Set("Content-Type", "application/json")

	//Alur: ambil claims dari context, populate id dengan context id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Failed Claims", http.StatusUnauthorized)
	}

	data, err := uh.Service.GetUserProfile(r.Context(), claims.Id)
	if err != nil {
		GenerateError(rw, "", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Update profile
// @description User update profile
// @tags User
// @accept json
// @Produce json
// @Param user body models.UserRequest true "Update Account Information"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /user [PUT]
// @Security BearerAuth
func (uh *UserHandler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var req = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
	}

	//Alur: ambil claims dari context, populate id dengan context id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Failed Claims", http.StatusUnauthorized)
	}

	//validate
	if err := utils.ValidateRequest(req.Name, req.Email, req.Password); len(err) > 0 {
		//Access setiap error, join ke joinedError, return sebagai message
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		GenerateError(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.UpdateUserProfile(r.Context(), claims.Id, req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Admin get users
// @description Return all existing users
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
	rw.Header().Set("Content-Type", "application/json")

	data, err := uh.Service.GetAllUsers(r.Context())
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		GenerateError(rw, "Gagal Encode", http.StatusInternalServerError)
		return
	}
}

// @Summary Admin get user profile
// @description get user profile by their unique ID
// @tags Admin
// @accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.AdminUserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /admin/user/{id} [GET]
// @Security BearerAuth
func (uh *UserHandler) AdminGetUserProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		GenerateError(rw, "Invalid ID", http.StatusBadRequest)
	}

	data, err := uh.Service.GetUserById(r.Context(), userID)
	if err != nil {
		GenerateError(rw, "", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// @Summary Admin delete user
// @description Admin removes user account by their unique ID
// @tags Admin
// @accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.AdminUserResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/user{id} [DELETE]
// @Security BearerAuth
func (uh *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		GenerateError(rw, "Invalid ID", http.StatusBadRequest)
	}

	response, err := uh.Service.DeleteUser(r.Context(), userID)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		GenerateError(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
