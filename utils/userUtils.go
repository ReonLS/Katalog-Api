package utils

import (
	"errors"
	"net/mail"
	"strings"
)

type Role string

const (
	RoleUser  Role = "User"
	RoleAdmin Role = "Admin"
)

func ValidateRequest(name string, email string, password string) []error {
	var listErr []error

	//validasi nama
	if name == "" || len(strings.TrimSpace(name)) == 0 {
		listErr = append(listErr, errors.New("Name must not be empty!"))
	} else if len(name) <= 2 {
		listErr = append(listErr, errors.New("Name must be at least 3 Characters long!"))
	}

	//validasi email
	if err := validateEmail(email); err != nil {
		for _, er := range err {
			listErr = append(listErr, er)
		}
	}

	//validasi password, nerima 8
	if err := validatePassword(password); err != nil {
		for _, er := range err {
			listErr = append(listErr, er)
		}
	}
	//berarti aman
	return listErr
}

func ValidateLogin(email string, password string) []error {
	var listErr []error

	//validasi email
	if err := validateEmail(email); err != nil {
		for _, er := range err {
			listErr = append(listErr, er)
		}
	}

	//validasi password, nerima 8
	if err := validatePassword(password); err != nil {
		for _, er := range err {
			listErr = append(listErr, er)
		}

	}
	//berarti aman
	return listErr
}

// mungkin perlu implement unique email
func validateEmail(email string) []error {
	var listErr []error

	if email == "" || len(strings.TrimSpace(email)) == 0 {
		listErr = append(listErr, errors.New("Email May Not Be Empty!"))
	} else if _, err := mail.ParseAddress(email); err != nil {
		listErr = append(listErr, errors.New("Invalid Email Format"))
	}

	return listErr
}

func validatePassword(Password string) []error {
	var listErr []error

	if Password == "" || len(strings.TrimSpace(Password)) == 0 {
		listErr = append(listErr, errors.New("Password must not be empty!"))
	} else if len(Password) <= 7 {
		listErr = append(listErr, errors.New("Password must be at least 8 Characters long!"))
	}

	return listErr
}
