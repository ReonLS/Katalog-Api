package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/mail"
	"simple-product-api/models"
	"simple-product-api/service"
	"simple-product-api/utils"
	"strings"
)

type UserHandler struct {
	Service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler{
	return &UserHandler{Service: service}
}

func ValidateRegister (req *models.UserRequest) (error){
	//validasi nama
	if req.Name == "" {
		return errors.New("Name must not be empty!")
	}
	if len(req.Name) < 3 {
		return errors.New("Name must be at least 3 Characters long!")
	}
	
	//validasi email
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	
	//validasi password, nerima 8
	if err := ValidatePassword(req.Password); err != nil{
		return err
	}
	//berarti aman
	return nil
}

func ValidateLogin (req *models.LoginRequest) (error){
	//validasi email
	if err := validateEmail(req.Email); err != nil {
		return err
	}
	
	//validasi password, nerima 8
	if err := ValidatePassword(req.Password); err != nil{
		return err
	}
	//berarti aman
	return nil
}

//mungkin perlu implement unique email
func validateEmail(email string) error {
	if email == ""{
		return errors.New("Email May Not Be Empty!")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return err
	}
	return nil
}

func ValidatePassword(Password string) (error){
	if Password == "" {
		return errors.New("Password must not be empty!")
	}
	if len(Password) <= 7 {
		return errors.New("Password must be at least 8 Characters long!")
	}
	return nil
}

func (uh *UserHandler) Register(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var request = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)

	//validate
	if err := ValidateRegister(request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "Err Request", http.StatusBadRequest)
		return
	}

	response, err := uh.Service.Register(r.Context(), request)
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

func (uh *UserHandler) Login(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var request = &models.LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	
	//validate
	if err := ValidateLogin(request); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "Error Request", http.StatusBadRequest)
		return
	}

	//could be either error or token
	token, err := uh.Service.Login(r.Context(), request)
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

func (uh *UserHandler) GetAllUsers(rw http.ResponseWriter, r *http.Request){
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

//User
func (uh *UserHandler) GetProfile(rw http.ResponseWriter, r *http.Request){
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

//Admin
func (uh *UserHandler) AdminGetUserProfile(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/admin/user/")

	data, err := uh.Service.GetUserById(r.Context(), idstring)
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

//user
func (uh *UserHandler) UpdateProfile(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//decode
	var request = &models.UserRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	
	//Alur: ambil claims dari context, populate id dengan context id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "Failed Claims", http.StatusUnauthorized)
	}

	response, err := uh.Service.UpdateUserProfile(r.Context(), claims.Id, request)
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

func (uh *UserHandler) DeleteUser(rw http.ResponseWriter, r *http.Request){
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/user/")
	
	response, err := uh.Service.DeleteUser(r.Context(), idstring)
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

