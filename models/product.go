package models

//repository pattern
type ProductRepository interface{
	GetProduct() ([]*Product, error)
	InsertProduct(req *Product) (*Product, error)
	UpdateProductByID(id int, req *Product) (*Product, error)
	DeleteProductByID(id int) (*Product, error)
}

type Product struct{
	Id int
	UserId int
	Namaprod string
	Kategori string
	Price float64
	Stock int
}

//for now, req ada user id jd bs ambil produk apa yg milik user id itu (next auto cek based by authorization on middleware)
//jadi cuman bs liat produk unik milik user
type ProductRequest struct{
	Namaprod string `json:"namaprod" binding:"required"`
	Kategori string `json:"kategori" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
	Stock int `json:"stock" binding:"required,gt=0"`
}

type ProductResponse struct{
	Id int `json:"id"`
	UserId int `json:"userid"`
	Namaprod string `json:"namaprod"`
	Kategori string `json:"kategori"`
	Price float64 `json:"price"`
	Stock int `json:"stock"`
}

