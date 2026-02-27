package service

import (
	"context"
	"simple-product-api/models"
	"github.com/google/uuid"
)

type ProductService struct{
	Repo models.ProductRepository
}


func NewProductService(repo models.ProductRepository) *ProductService{
	return &ProductService{Repo: repo}
}

func (pr *ProductService) ToProductResponse(p *models.Product) *models.UserProductResponse{
	//transform dari domain struct(db) jd response (json-embedded)
	return &models.UserProductResponse{
		Id: p.Id,
		Namaprod: p.Namaprod,
		Kategori: p.Kategori,
		Price: p.Price,
		Stock: p.Stock,
	}
}

func (pr *ProductService) ToAdminProductResponse(p *models.Product) *models.AdminProductResponse{
	//transform dari domain struct(db) jd response (json-embedded)
	return &models.AdminProductResponse{
		Id: p.Id,
		UserId: p.UserId,
		Namaprod: p.Namaprod,
		Kategori: p.Kategori,
		Price: p.Price,
		Stock: p.Stock,
	}
}

func (pr *ProductService) GetProductByUserID(ctx context.Context, id string)([]*models.UserProductResponse, error){
	//Alur : Nerima domain struct, transform jadi response 
	var dataResp []*models.UserProductResponse

	data, err := pr.Repo.GetProductByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	//for loop access masing2
	for _, rows := range data{
		dataResp = append(dataResp, pr.ToProductResponse(rows))
	}

	//aman berarti
	return dataResp, nil
}

func (pr *ProductService) InsertProduct(ctx context.Context, userid string, req *models.ProductRequest) (*models.UserProductResponse, error){
	//Alur : Nerima domain struct, generate product.response
	var data = &models.Product{
		Id: uuid.New().String(),
		UserId: userid,
		Namaprod: req.Namaprod,
		Kategori: req.Kategori,
		Price: req.Price,
		Stock: req.Stock,
	}

	product, err := pr.Repo.InsertProduct(ctx, userid, data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) UpdateProductByID(ctx context.Context, id string, req *models.ProductRequest) (*models.UserProductResponse, error){
	//Alur : Nerima domain struct, generate product.response
	var data = &models.Product{
		Id: uuid.New().String(),
		Namaprod: req.Namaprod,
		Kategori: req.Kategori,
		Price: req.Price,
		Stock: req.Stock,
	}

	product, err := pr.Repo.UpdateProductByID(ctx, id, data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) DeleteProductByID(ctx context.Context, id string) (*models.UserProductResponse, error){
	//Alur : Nerima domain struct, generate product.response

	product, err := pr.Repo.DeleteProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	//aman
	return pr.ToProductResponse(product), nil
}