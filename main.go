// @title Simple Product Api Visualization
// @version 1.0
// @description A Product Catalogue API with Swagger Documentation

// @host localhost:8080
// @Basepath /

// @security.Definitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT Token
package main

import (
	"fmt"
	"net/http"
	"simple-product-api/config"
	"simple-product-api/handler"
	"simple-product-api/repository"
	"simple-product-api/route"
	"simple-product-api/service"
	_ "simple-product-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/go-chi/chi/v5"
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
	route.RouteSetup(r)

	r.Get("/swagger/*", httpSwagger.WrapHandler)
	err := http.ListenAndServe(config.GetMainPort(), r)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
