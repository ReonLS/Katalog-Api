package route

import (
	"simple-product-api/handler"
	"simple-product-api/middleware"
	"simple-product-api/utils"

	"github.com/go-chi/chi/v5"
)

//ngebuat instans of product handler
type Route struct{
	ProdHandler handler.ProductHandler
	UserHandler handler.UserHandler
}

func NewProductRoute(product *handler.ProductHandler, user *handler.UserHandler) *Route{
	return &Route{ProdHandler: *product, UserHandler: *user}
}

//centralized handler func for /product
func (route *Route) RouteSetup (r chi.Router){
	//Public
	r.Post("/register", route.UserHandler.Register)
	r.Post("/login", route.UserHandler.Login)

	//Endpoint: /user
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticateJWT)
		r.Get("/user", route.UserHandler.GetProfile)
		r.Put("/user", route.UserHandler.UpdateProfile) 
	})

	//Endpoint : /admin/user
	r.Route("/admin", func(r chi.Router) {
		r.Use(middleware.AuthenticateJWT)
		r.Use(middleware.AuthenticateRole(utils.RoleAdmin))

		r.Get("/user", route.UserHandler.GetAllUsers) //admin
		r.Get("/user/{id}", route.UserHandler.AdminGetUserProfile) //admin
		r.Delete("/user/{id}", route.UserHandler.DeleteUser) //admin
	})

	//Endpoint: /product
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthenticateJWT)
		r.Get("/product", route.ProdHandler.GetProduct)
		r.Post("/product", route.ProdHandler.InsertProduct)
		r.Put("/product/{id}", route.ProdHandler.UpdateProductByID)
		r.Delete("/product/{id}", route.ProdHandler.DeleteProductByID)
	})
}