package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"simple-product-api/models"
	"simple-product-api/service"
	"strconv"
	"strings"
)

// buat instans DB, jdi layer handler bisa exec query
// sebenarnya ini layer repo, yg diakses service, yg diakses handler tp for now okela
type ProductHandler struct {
	Service service.ProductService
}

func NewProductHandler(service *service.ProductService) *ProductHandler{
	return &ProductHandler{Service: *service}
}

// Fungsi utama responds to a request (getting all product)
func (ph *ProductHandler) GetProduct(rw http.ResponseWriter, r *http.Request) {
	//Alur : Nerima response, encode jadi json 
	rw.Header().Set("Content-Type", "application/json")

	products, err := ph.Service.GetProduct()
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
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
	fmt.Println("Masuk POST")

	//membuat tampungan
	var request = &models.ProductRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		//dianggap client salah kirim input
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	//logikanya gagal kebentuk, berarti user kirim faulty request
	response, err := ph.Service.InsertProduct(request)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
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
	
	//parsing id dari path
	path := r.URL.Path //{/{id}}
	stringId := strings.TrimPrefix(path, "/product/") //{id}
	fmt.Println("Masuk PUT", stringId)

	id, err := strconv.Atoi(stringId)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	//tampungan decode
	var request = &models.ProductRequest{}

	err = json.NewDecoder(r.Body).Decode(&request)
	defer r.Body.Close()

	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	//panggil service func
	response, err := ph.Service.UpdateProductByID(id, request)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//encode update untuk write ke stream
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
}

func (ph *ProductHandler) DeleteProductByID(rw http.ResponseWriter, r *http.Request) {
	//alur : set header -> ambil ID dari url, decode, jalankan query, encode, response
	rw.Header().Set("Content-Type", "application/json")

	//generate id from path
	path := r.URL.Path
	idstring := strings.TrimPrefix(path, "/product/")
	id, err := strconv.Atoi(idstring)
	fmt.Println("Masuk DELETE", id)

	if err != nil{
		http.Error(rw, "", http.StatusBadRequest)
		return
	}

	//jalankan query
	response, err := ph.Service.DeleteProductByID(id)
	if err != nil {
		http.Error(rw, "", http.StatusBadRequest)
		return
	}
	//berarti aman
	rw.WriteHeader(http.StatusOK)

	//tembak ke stream
	err = json.NewEncoder(rw).Encode(response)
	if err != nil {
		http.Error(rw, "", http.StatusInternalServerError)
		return
	}
}
