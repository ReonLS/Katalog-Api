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
	prodRepo := repository.NewProductRepo(db)
	prodService := service.NewProductService(prodRepo)
	prodHandler := handler.NewProductHandler(prodService)

	//user
	userRepo := repository.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	route := route.NewProductRoute(prodHandler, userHandler)

	mux := http.NewServeMux()
	route.Product(mux)
	route.User(mux)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
