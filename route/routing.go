package route

import "net/http"
import "simple-product-api/handler"

//ngebuat instans of product handler
type Route struct{
	ProdHandler *handler.ProductHandler
	UserHandler *handler.UserHandler
}

//centralized handler func for /product
func (r *Route) Product (mux *http.ServeMux){
	mux.HandleFunc("GET /product", r.ProdHandler.GetProduct)
	mux.HandleFunc("POST /product", r.ProdHandler.InsertProduct)
	mux.HandleFunc("PUT /product/{id}", r.ProdHandler.UpdateProductByID)
	mux.HandleFunc("DELETE /product/{id}", r.ProdHandler.DeleteProductByID)
}

func (r *Route) User (mux *http.ServeMux){
	mux.HandleFunc("GET /user", r.UserHandler.GetAllUsers)
	mux.HandleFunc("GET /user/{id}", r.UserHandler.GetUserbyId)
	mux.HandleFunc("POST /user", r.UserHandler.CreateUser)
	mux.HandleFunc("PUT /user/{id}", r.UserHandler.UpdateUser)
	mux.HandleFunc("DELETE /user/{id}", r.UserHandler.DeleteUser)
}
