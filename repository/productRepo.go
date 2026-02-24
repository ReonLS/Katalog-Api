package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"simple-product-api/models"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo{
	return &ProductRepo{DB: db}
}

func (pr *ProductRepo) GetProduct() ([]*models.Product, error) {
	//Alur : Generate query, return domain struct

	var data []*models.Product

	rows, err := pr.DB.Query("Select * from product")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rowData = &models.Product{}

		if err = rows.Scan(&rowData.Id, &rowData.Namaprod, &rowData.Kategori, &rowData.Price, &rowData.Stock); err != nil {
			return nil, err
		}
		data = append(data, rowData)
	}
	//semua aman
	return data, nil
}

func (pr *ProductRepo) InsertProduct(prod *models.Product) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (ngamnbil id dari hasil auto increment table)

	//query row exec query dan return data including ID pake returning, sebagai return value
	query := "Insert into product (namaprod, kategori, price, stock) values (?,?,?,?)"

	//logikanya tu karna domain struct punya value sama aja, disini hasil queryrow return ID
	//karna cukup butuh mappingan last inserted id untuk generate id product baru
	result, err := pr.DB.Exec(query, prod.Namaprod, prod.Kategori, prod.Price, prod.Stock)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, errors.New("Product Not Created")
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	prod.Id = int(id)

	//artinya aman
	return prod, nil
}

func (pr *ProductRepo) UpdateProductByID(id int, product *models.Product) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (semua info property udh dari request)

	query := "update product set namaprod=?, kategori=?, price=?, stock=? where id = ?"
	res, err := pr.DB.Exec(query, product.Namaprod, product.Kategori, product.Price, product.Stock, id)
	if err != nil {
		return nil, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rows == 0 {
		return nil, err
	}

	//berarti aman
	return product, nil
}

// in proper api, query ttp delete unique product id, tp middleware yg bakal authenticate user
// untuk ensure product ini milik currentuserloginid
func (pr *ProductRepo) DeleteProductByID(id int) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (ngamnbil id dari hasil auto increment table)

	var product = &models.Product{
		Id: id,
	}

	//query select based id, untuk dpt info deleted baru jalanin delete query
	err := pr.DB.QueryRow("select namaprod,kategori,price,stock from product where id = ?", id).
		Scan(&product.Namaprod, &product.Kategori, &product.Price, &product.Stock)

	fmt.Println(product.Id, product.Namaprod, product.Kategori, product.Price, product.Stock)
	if err != nil {
		return nil, err
	}

	result, err := pr.DB.Exec("Delete from product where id = ?", id)
	if err != nil {
		return nil, err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAff == 0 {
		return nil, errors.New("Product Not Found!")
	}
	//artinya aman
	return product, nil
}
