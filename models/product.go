package models

import (
	"context"
	"simple-product-api/utils"
)

// repository pattern
type ProductRepository interface {
	GetProductByUserID(ctx context.Context, userID string) ([]*Product, error) //id dari context
	GetProductByProdID(ctx context.Context, prodID string) (*Product, error)   //id dari context
	InsertProduct(ctx context.Context, userID string, req *Product) (*Product, error)
	UpdateProductByID(ctx context.Context, id string, req *Product) (*Product, error)
	DeleteProductByID(ctx context.Context, id string) (*Product, error)
}

type Product struct {
	Id       string
	UserId   string
	Namaprod string
	Kategori utils.Category
	Price    float64
	Stock    int
}

// for now, req ada user id jd bs ambil produk apa yg milik user id itu (next auto cek based by authorization on middleware)
// jadi cuman bs liat produk unik milik user
type ProductRequest struct {
	Namaprod string         `json:"namaprod" example:"Kemeja Putih"`
	Kategori utils.Category `json:"kategori" example:"Baju"`
	Price    float64        `json:"price" example:"120.5"`
	Stock    int            `json:"stock" example:"15"`
}

type UserProductResponse struct {
	Id       string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Namaprod string         `json:"namaprod" example:"Kemeja Putih"`
	Kategori utils.Category `json:"kategori" example:"Baju"`
	Price    float64        `json:"price" example:"120.5"`
	Stock    int            `json:"stock" example:"15"`
}

type AdminProductResponse struct {
	Id       string         `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserId   string         `json:"userid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Namaprod string         `json:"namaprod" example:"Kemeja Putih"`
	Kategori utils.Category `json:"kategori" example:"Baju"`
	Price    float64        `json:"price" example:"120.5"`
	Stock    int            `json:"stock" example:"15"`
}
