// @title Product Katalog
// @version 1.0
// @description A RESTful product catalogue API built in Go, featuring JWT-based authentication (HS256) and role-based access control. Protected endpoints require a Bearer token obtained via `POST /api/v1/login`. Two roles are supported: `User` and `Admin`, each with distinct access levels.
// @host localhost:8080
// @Basepath /api/v1
// @tag.name Public
// @tag.description Publicly accessible endpoints
// @tag.name User
// @tag.description Operations restricted only to user with role "User"
// @tag.name Admin
// @tag.description Operations restricted only to user with role "Admin"
// @SecurityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT Token
package main

import (
	"fmt"
	"net/http"
	"simple-product-api/config"
	_ "simple-product-api/docs"
	"simple-product-api/handler"
	"simple-product-api/repository"
	"simple-product-api/route"
	"simple-product-api/service"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	//init DB & Env
	config.Init()
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
	r.Route("/api/v1", func(r chi.Router) {
		route.RouteSetup(r)
	})
	
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	err := http.ListenAndServe(config.GetMainPort(), r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
