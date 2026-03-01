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

// Fungsi utama responds to a request (getting all product)
func (ph *ProductHandler) GetProduct(rw http.ResponseWriter, r *http.Request) {
	//Alur : Nerima response, encode jadi json
	rw.Header().Set("Content-Type", "application/json")

	//placeholder ngambils claims dari context, ambil id
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "No Information", http.StatusUnauthorized)
	}

	products, err := ph.Service.GetUserProduct(r.Context(), claims.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//Best Approach, more memory efficient
	err = json.NewEncoder(rw).Encode(products)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *ProductHandler) AdminGetProductUser(rw http.ResponseWriter, r *http.Request) {
	//Alur : Nerima response, encode jadi json
	rw.Header().Set("Content-Type", "application/json")

		//Parsing id form path, validation
	userID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(userID); err != nil{
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	products, err := ph.Service.AdminGetUserProduct(r.Context(), userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//Best Approach, more memory efficient
	err = json.NewEncoder(rw).Encode(products)
	if err != nil {
		//server-side error
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *ProductHandler) InsertProduct(rw http.ResponseWriter, r *http.Request) {
	//Alur real life : nerima json -> decode dan simpan di tampungan, exec query, generate respon
	rw.Header().Set("Content-Type", "application/json")

	//membuat tampungan
	var req = &models.ProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		//dianggap client salah kirim input
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//validasi input
	if err := utils.ValidateProduct(req.Namaprod, string(req.Kategori), req.Price, req.Stock); len(err) > 0 {
		//Access setiap error, join ke joinedError, return sebagai message
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		http.Error(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	//ambil userid from context
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "No Information", http.StatusUnauthorized)
		return
	}

	//logikanya gagal kebentuk, berarti user kirim faulty request
	response, err := ph.Service.InsertProduct(r.Context(), claims.Id, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//return status http object berhasil dibentuk
	rw.WriteHeader(http.StatusCreated)

	//tampilin di endpoint sbg response request client
	json.NewEncoder(rw).Encode(response)
}

func (ph *ProductHandler) UpdateProductByID(rw http.ResponseWriter, r *http.Request) {
	//alur : set header, take id from url.path, decode req.body, call service func, generate respons
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	prodID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(prodID); err != nil{
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	//tampungan decode
	var req = &models.ProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	//validasi input
	if err := utils.ValidateProduct(req.Namaprod, string(req.Kategori), req.Price, req.Stock); len(err) > 0 {
		//Access setiap error, join ke joinedError, return sebagai message
		var joinedError []string
		for _, each := range err {
			joinedError = append(joinedError, each.Error())
		}

		http.Error(rw, strings.Join(joinedError, "\n"), http.StatusBadRequest)
		return
	}

	//ambil userId from Claims
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "Need Authorization", http.StatusUnauthorized)
		return
	}

	//panggil service func
	response, err := ph.Service.UpdateProductByID(r.Context(), prodID, claims.Id, req)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//encode update untuk write ke stream
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (ph *ProductHandler) DeleteProductByID(rw http.ResponseWriter, r *http.Request) {
	//alur : set header -> ambil ID dari url, decode, jalankan query, encode, response
	rw.Header().Set("Content-Type", "application/json")

	//Parsing id form path, validation
	prodID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(prodID); err != nil{
		http.Error(rw, "Invalid ID", http.StatusBadRequest)
	}

	//ambil userId from Claims
	claims, ok := utils.GetClaimsFromContext(r.Context())
	if !ok {
		http.Error(rw, "Need Authorization", http.StatusUnauthorized)
		return
	}

	//jalankan query
	response, err := ph.Service.DeleteProductByID(r.Context(), prodID, claims.Id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//tembak ke stream
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
}
