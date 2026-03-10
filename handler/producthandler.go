package handler

import (
	"encoding/json"
	"net/http"
	"simple-product-api/models"
	"simple-product-api/service"
	"simple-product-api/utils"
	"strings"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// buat instans DB, jdi layer handler bisa exec query
// sebenarnya ini layer repo, yg diakses service, yg diakses handler tp for now okela
type ProductHandler struct {
	Service *service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler {
	return &ProductHandler{Service: service}
}

// @Summary Get all product
// @description Retrieve all user's product automatically by userID from context
// @tags User
// @accept json
// @Produce json
// @Success 200 {array} models.UserProductResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /user/product [GET]
// @Security BearerAuth
func (ph *ProductHandler) GetProduct(rw http.ResponseWriter, r *http.Request) {
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "No Information", http.StatusUnauthorized)
		return
	}

	products, err := ph.Service.GetUserProduct(r.Context(), claims.Id)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, products)
}

// @Summary Insert product
// @description Returns inserted product's information
// @tags User
// @accept json
// @Produce json
// @Param product body models.ProductRequest true "Product credentials"
// @Success 201 {array} models.UserProductResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /user/product [POST]
// @Security BearerAuth
func (ph *ProductHandler) InsertProduct(rw http.ResponseWriter, r *http.Request) {
	var req = &models.ProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	
	if err := utils.ValidateProduct(req.Namaprod, string(req.Kategori), req.Price, req.Stock); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}
		GenerateError(rw, strings.Join(joinedError, " ,"), http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "No Information", http.StatusUnauthorized)
		return
	}

	response, err := ph.Service.InsertProduct(r.Context(), claims.Id, req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusCreated, response)
}

// @Summary Update product by ID
// @description Returns updated product's information
// @tags User
// @accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param product body models.ProductRequest true "Update Product"
// @Success 200 {array} models.UserProductResponse
// @Failure 400 {object} models.BadRequestResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /user/product/{id} [PUT]
// @Security BearerAuth
func (ph *ProductHandler) UpdateProductByID(rw http.ResponseWriter, r *http.Request) {
	prodID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(prodID); err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var req = &models.ProductRequest{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := utils.ValidateProduct(req.Namaprod, string(req.Kategori), req.Price, req.Stock); len(err) > 0 {
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}
		GenerateError(rw, strings.Join(joinedError, " ,"), http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Need Authorization", http.StatusUnauthorized)
		return
	}

	response, err := ph.Service.UpdateProductByID(r.Context(), prodID, claims.Id, req)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Delete product by ID
// @description Returns deleted product's information
// @tags User
// @accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {array} models.UserProductResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse 
// @Router /user/product/{id} [DELETE]
// @Security BearerAuth
func (ph *ProductHandler) DeleteProductByID(rw http.ResponseWriter, r *http.Request) {
	prodID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(prodID); err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		GenerateError(rw, "Need Authorization", http.StatusUnauthorized)
		return
	}

	response, err := ph.Service.DeleteProductByID(r.Context(), prodID, claims.Id)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Get user product by UserID
// @description Returns all product belonging to a user
// @tags Admin
// @accept json
// @Produce json
// @param id path string true "User ID"
// @Success 200 {array} models.UserProductResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/{id}/product [GET]
// @Security BearerAuth
func (ph *ProductHandler) AdminGetProductUser(rw http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := ph.Service.GetUserProduct(r.Context(), userID)
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}

// @Summary Get all product
// @description Returns all existing product
// @tags Admin
// @accept json
// @Produce json
// @Success 200 {array} models.AdminProductResponse
// @Failure 401 {object} models.UnauthorizedResponse 
// @Failure 403 {object} models.ForbiddenResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /admin/product [GET]
// @Security BearerAuth
func (ph *ProductHandler) AdminGetAllProduct(rw http.ResponseWriter, r *http.Request) {
	response, err := ph.Service.AdminGetAllProduct(r.Context())
	if err != nil {
		GenerateError(rw, err.Error(), http.StatusBadRequest)
		return
	}
	WriteJSON(rw, http.StatusOK, response)
}
