package repository

import (
	"context"
	"database/sql"
	"errors"
	"simple-product-api/models"
)

type ProductRepo struct {
	DB *sql.DB
}

func NewProductRepo(db *sql.DB) *ProductRepo {
	return &ProductRepo{DB: db}
}

func (pr *ProductRepo) GetProductByUserID(ctx context.Context, userID string) ([]*models.Product, error) {
	//Alur : Generate query, return domain struct

	var data []*models.Product

	rows, err := pr.DB.QueryContext(ctx, "Select * from product where userid = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var rowData = &models.Product{}

		if err = rows.Scan(&rowData.Id, &rowData.UserId, &rowData.Namaprod, &rowData.Kategori, &rowData.Price, &rowData.Stock); err != nil {
			return nil, err
		}
		data = append(data, rowData)
	}
	//semua aman
	return data, nil
}

func (pr *ProductRepo) GetProductByProdID(ctx context.Context, prodID string) (*models.Product, error) {
	//Alur : Generate query, return domain struct

	var data = &models.Product{}

	res := pr.DB.QueryRowContext(ctx, "Select * from product where id = ?", prodID)
	if err := res.Err(); err != nil {
		return nil, err
	}

	err := res.Scan(&data.Id, &data.UserId, &data.Namaprod, &data.Kategori, &data.Price, &data.Stock)
	if err != nil {
		return nil, err
	}
	//semua aman
	return data, nil
}

func (pr *ProductRepo) InsertProduct(ctx context.Context, userid string, prod *models.Product) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (ngamnbil id dari hasil auto increment table)

	//query row exec query dan return data including ID pake returning, sebagai return value
	query := "Insert into product (id, userid, namaprod, kategori, price, stock) values (?,?,?,?,?,?)"

	//logikanya tu karna domain struct punya value sama aja, disini hasil queryrow return ID
	//karna cukup butuh mappingan last inserted id untuk generate id product baru
	result, err := pr.DB.ExecContext(ctx, query, prod.Id, prod.UserId, prod.Namaprod, prod.Kategori, prod.Price, prod.Stock)
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
	//artinya aman
	return prod, nil
}

func (pr *ProductRepo) UpdateProductByID(ctx context.Context, prodID string, product *models.Product) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (semua info property udh dari request)

	query := "update product set namaprod=?, kategori=?, price=?, stock=? where id = ?"
	res, err := pr.DB.ExecContext(ctx, query, product.Namaprod, product.Kategori, product.Price, product.Stock, prodID)
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
func (pr *ProductRepo) DeleteProductByID(ctx context.Context, id string) (*models.Product, error) {
	//Alur : Jalanin query, return domain struct (ngamnbil id dari hasil auto increment table)

	var product = &models.Product{}
	//query select based id, untuk dpt info deleted baru jalanin delete query
	err := pr.DB.QueryRowContext(ctx, "select * from product where id = ?", id).
		Scan(&product.Id, &product.UserId, &product.Namaprod, &product.Kategori, &product.Price, &product.Stock)

	if err != nil {
		return nil, err
	}

	result, err := pr.DB.ExecContext(ctx, "Delete from product where id = ?", id)
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
