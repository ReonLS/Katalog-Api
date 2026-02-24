package service

import "simple-product-api/models"

type ProductService struct{
	Repo models.ProductRepository
}


func NewProductService(repo models.ProductRepository) *ProductService{
	return &ProductService{Repo: repo}
}

func (pr *ProductService) ToProductResponse(p *models.Product) *models.ProductResponse{
	//transform dari domain struct(db) jd response (json-embedded)
	return &models.ProductResponse{
		Id: p.Id,
		Namaprod: p.Namaprod,
		Kategori: p.Kategori,
		Price: p.Price,
		Stock: p.Stock,
	}
}

func (pr *ProductService) GetProduct()([]*models.ProductResponse, error){
	//Alur : Nerima domain struct, transform jadi response 
	var dataResp []*models.ProductResponse

	data, err := pr.Repo.GetProduct()
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

func (pr *ProductService) InsertProduct(req *models.ProductRequest) (*models.ProductResponse, error){
	//Alur : Nerima domain struct, generate product.response
	var data = &models.Product{
		Namaprod: req.Namaprod,
		Kategori: req.Kategori,
		Price: req.Price,
		Stock: req.Stock,
	}

	product, err := pr.Repo.InsertProduct(data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) UpdateProductByID(id int, req *models.ProductRequest) (*models.ProductResponse, error){
	//Alur : Nerima domain struct, generate product.response
	var data = &models.Product{
		Id: id,
		Namaprod: req.Namaprod,
		Kategori: req.Kategori,
		Price: req.Price,
		Stock: req.Stock,
	}

	product, err := pr.Repo.UpdateProductByID(id, data)
	if err != nil {
		return nil, err
	}

	//aman
	return pr.ToProductResponse(product), nil
}

func (pr *ProductService) DeleteProductByID(id int) (*models.ProductResponse, error){
	//Alur : Nerima domain struct, generate product.response

	product, err := pr.Repo.DeleteProductByID(id)
	if err != nil {
		return nil, err
	}
	//aman
	return pr.ToProductResponse(product), nil
}