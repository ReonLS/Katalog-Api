package main

import (
	"fmt"
	"net/http"
	"simple-product-api/config"
	"simple-product-api/handler"
	"simple-product-api/route"
	"simple-product-api/repository"
	"simple-product-api/service"
)

func main() {
	//init DB
	db := config.Connect()
	defer db.Close()

	//products
	prodRepo := repository.ProductRepo{DB: db}
	prodService := service.ProductService{Repo: &prodRepo}
	prodHandler := handler.ProductHandler{Service: &prodService}
	prodRoute := route.Route{ProdHandler: &prodHandler}

	//user
	userRepo := repository.UserRepo{DB: db}
	userService := service.UserService{Repo: &userRepo}
	userHandler := handler.UserHandler{Service: &userService}
	userRoute := route.Route{UserHandler: &userHandler}

	mux := http.NewServeMux()
	prodRoute.Product(mux)
	userRoute.User(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
