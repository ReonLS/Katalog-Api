package handler

import (
	"encoding/json"
	"net/http"
	"simple-product-api/models"
	"simple-product-api/service"
	"strconv"
	"strings"
	"net/mail"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{Service: service}
}

//mungkin perlu implement unique email, harusnya ini di models
func (uh *UserHandler) validateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func (uh *UserHandler) GetAllUsers(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	data, err := uh.Service.GetAllUsers()
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

func (uh *UserHandler) GetUserbyId(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/user/")
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	data, err := uh.Service.GetUserbyId(id)
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

func (uh *UserHandler) CreateUser(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var request = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	//validate email
	if err := uh.validateEmail(request.Email); err != nil {
		http.Error(rw, "Email not valid", http.StatusBadRequest)
		return
	}

	response, err := uh.Service.CreateUser(request)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
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

func (uh *UserHandler) UpdateUser(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/user/")
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	//decode
	var request = &models.UserRequest{}
	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
	}

	//validate email
	if err := uh.validateEmail(request.Email); err != nil {
		http.Error(rw, "Email not valid", http.StatusBadRequest)
		return
	}

	response, err := uh.Service.UpdateUser(id, request)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
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

func (uh *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/user/")
	id, err := strconv.Atoi(idstring)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}

	response, err := uh.Service.DeleteUser(id)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
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

