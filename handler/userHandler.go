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

		http.Error(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.Register(r.Context(), req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) Login(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var req = &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "Error Request", http.StatusBadRequest)
		return
	}

	//validate
	if err := utils.ValidateLogin(req.Email, req.Password); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		http.Error(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	//could be either error or token
	token, err := uh.Service.Login(r.Context(), req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(token)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) GetAllUsers(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	data, err := uh.Service.GetAllUsers(r.Context())
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		http.Error(rw, "Gagal Encode", http.StatusInternalServerError)
		return
	}
}

// User
func (uh *UserHandler) GetProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//Alur: ambil claims dari context, populate id dengan context id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "Failed Claims", http.StatusUnauthorized)
	}

	data, err := uh.Service.GetUserProfile(r.Context(), claims.Id)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Admin
func (uh *UserHandler) AdminGetUserProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	data, err := uh.Service.GetUserById(r.Context(), userID)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(data)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

// user
func (uh *UserHandler) UpdateProfile(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var req = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	//Alur: ambil claims dari context, populate id dengan context id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "Failed Claims", http.StatusUnauthorized)
	}

	//validate
	if err := utils.ValidateRequest(req.Name, req.Email, req.Password); len(err) > 0 {
		//Access setiap error, join ke joinedError, return sebagai message
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		http.Error(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	response, err := uh.Service.UpdateUserProfile(r.Context(), claims.Id, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (uh *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	response, err := uh.Service.DeleteUser(r.Context(), userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.WriteHeader(http.StatusOK)

	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
