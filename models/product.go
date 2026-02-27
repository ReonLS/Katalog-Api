package models

import "context"

//repository pattern
type ProductRepository interface{
	GetProductByUserID(ctx context.Context, id string) ([]*Product, error) //id dari context
	InsertProduct(ctx context.Context ,userid string, req *Product) (*Product, error)
	UpdateProductByID(ctx context.Context, id string, req *Product) (*Product, error)
	DeleteProductByID(ctx context.Context, id string) (*Product, error)
}

type Product struct{
	Id string
	UserId string
	Namaprod string
	Kategori string
	Price float64
	Stock int
}

//for now, req ada user id jd bs ambil produk apa yg milik user id itu (next auto cek based by authorization on middleware)
//jadi cuman bs liat produk unik milik user
type ProductRequest struct{
	Namaprod string `json:"namaprod"`
	Kategori string `json:"kategori"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
}

type UserProductResponse struct{
	Id string `json:"id"`
	Namaprod string `json:"namaprod"`
	Kategori string `json:"kategori"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
}

type AdminProductResponse struct{
	Id string `json:"id"`
	UserId string `json:"userid"`
	Namaprod string `json:"namaprod"`
	Kategori string `json:"kategori"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
}

