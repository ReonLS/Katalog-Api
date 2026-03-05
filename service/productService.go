package service

import (
	"context"
	"errors"
	"simple-product-api/models"

	"github.com/google/uuid"
)

type ProductService struct {
	Repo models.ProductRepository
}

func NewProductService(repo models.ProductRepository) *ProductService {
	return &ProductService{Repo: repo}
}

func (pr *ProductService) ToProductResponse(p *models.Product) *models.UserProductResponse {
	//transform dari domain struct(db) jd response (json-embedded)
	return &models.UserProductResponse{
		Id:       p.Id,
		Namaprod: p.Namaprod,
		Kategori: p.Kategori,
		Price:    p.Price,
		Stock:    p.Stock,
	}
}

func (pr *ProductService) ToAdminProductResponse(p *models.Product) *models.AdminProductResponse {
	//transform dari domain struct(db) jd response (json-embedded)
	return &models.AdminProductResponse{
		Id:       p.Id,
		UserId:   p.UserId,
		Namaprod: p.Namaprod,
		Kategori: p.Kategori,
		Price:    p.Price,
		Stock:    p.Stock,
	}
}

func (pr *ProductService) GetUserProduct(ctx context.Context, userID string) ([]*models.UserProductResponse, error) {
	//Alur : Nerima domain struct, transform jadi response
	var dataResp []*models.UserProductResponse

	data, err := pr.Repo.GetProductByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	//for loop access masing2
	for _, rows := range data {
		dataResp = append(dataResp, pr.ToProductResponse(rows))
	}

	//aman berarti
	return dataResp, nil
}

func (pr *ProductService) AdminGetUserProduct(ctx context.Context, id string) ([]*models.AdminProductResponse, error) {
	//Alur : Nerima domain struct, transform jadi response
	var dataResp []*models.AdminProductResponse

	data, err := pr.Repo.GetProductByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	//for loop access masing2
	for _, rows := range data {
		dataResp = append(dataResp, pr.ToAdminProductResponse(rows))
	}

	//aman berarti
	return dataResp, nil
}

func (pr *ProductService) InsertProduct(ctx context.Context, userid string, req *models.ProductRequest) (*models.UserProductResponse, error) {
	//Alur : Nerima domain struct, generate product.response
	var data = &models.Product{
		Id:       uuid.New().String(),
		UserId:   userid,
		Namaprod: req.Namaprod,
		Kategori: req.Kategori,
		Price:    req.Price,
		Stock:    req.Stock,
	}

	product, err := pr.Repo.InsertProduct(ctx, userid, data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) UpdateProductByID(ctx context.Context, ProdID string, userID string, req *models.ProductRequest) (*models.UserProductResponse, error) {
	//Alur : Nerima domain struct, generate product.response

	//compare userID.ctx & userID.repo
	data, err := pr.Repo.GetProductByProdID(ctx, ProdID)
	if err != nil {
		return nil, err
	}

	//validasi ownership
	if userID != data.UserId {
		//http forbidden
		return nil, errors.New("Must be your own products")
	}

	//replace dgn req's info
	data.Id = ProdID
	data.Namaprod = req.Namaprod
	data.Kategori = req.Kategori
	data.Price = req.Price
	data.Stock = req.Stock

	product, err := pr.Repo.UpdateProductByID(ctx, ProdID, data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) DeleteProductByID(ctx context.Context, prodID string, userID string) (*models.UserProductResponse, error) {
	//Alur : Nerima domain struct, generate product.response
	data, err := pr.Repo.GetProductByProdID(ctx, prodID)
	if err != nil {
		return nil, err
	}

	//validasi ownership
	if userID != data.UserId {
		//http forbidden
		return nil, errors.New("Must be your own products")
	}

	data, err = pr.Repo.DeleteProductByID(ctx, prodID)
	if err != nil {
		return nil, err
	}
	//aman
	return pr.ToProductResponse(data), nil
}
