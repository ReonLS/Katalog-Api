package main

import (
	"fmt"
	"net/http"
	"simple-product-api/config"
	"simple-product-api/handler"
	"simple-product-api/route"
	"simple-product-api/repository"
	"simple-product-api/service"
	"github.com/go-chi/chi/v5"
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

	r := chi.NewRouter()
	route.RouteSetup(r)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
