package services

import (
	"beego-test/models"
	"beego-test/repositories"
	"encoding/json"
	"fmt"
	"time"

	"github.com/astaxie/beego/orm"
)

//IProductService interface to access product API
type IProductService interface {
	//List of products in page (starts at 1), and rpp (rows per page)
	List(page int, rpp int) ([]models.Product, *APIError)

	//Find a product by its ID
	Find(id int64) (models.Product, *APIError)

	//Create a new product
	Create([]byte) (models.Product, *APIError)
}

type productService struct {
	repo repositories.IProductRepository
}

//NewProductService factory method to instantiate product service
func NewProductService(repo repositories.IProductRepository) IProductService {
	return &productService{repo: repo}
}

func (svc *productService) List(page int, rpp int) ([]models.Product, *APIError) {
	if page < 1 {
		return nil, &APIError{Code: 400, Message: "Page must be greater than zero"}
	}

	if rpp < 10 {
		return nil, &APIError{Code: 400, Message: "Minimum rows per page is 10"}
	}

	prod, err := svc.repo.List(page, rpp)
	if err != nil {
		return nil, &APIError{Code: 500, Message: fmt.Sprintf("Error getting product list: %s", err.Error())}
	}

	return prod, nil
}

func (svc *productService) Find(id int64) (models.Product, *APIError) {
	if id <= 0 {
		return models.Product{}, &APIError{Code: 400, Message: "product ID cannot be empty"}
	}

	prod, err := svc.repo.Find(id)
	if err == orm.ErrNoRows {
		return models.Product{}, &APIError{Code: 404, Message: fmt.Sprintf("Product with ID: %d is not found", id)}
	}

	if err != nil {
		return models.Product{}, &APIError{
			Code:    500,
			Message: fmt.Sprintf("Error getting product with ID: %d, %s", id, err.Error()),
		}
	}

	return prod, nil
}

func (svc *productService) Create(body []byte) (models.Product, *APIError) {
	//1. Parse request body
	prod := models.Product{}
	err := json.Unmarshal(body, &prod)
	if err != nil {
		return prod, &APIError{Code: 400, Message: "Invalid JSON request"}
	}

	//2. Validate product data
	switch {
	case prod.Name == "" || len(prod.Name) <= 0:
		return prod, &APIError{Code: 400, Message: "Product name is required"}
	case prod.Category == "" || len(prod.Category) <= 0:
		return prod, &APIError{Code: 400, Message: "Product name is required"}
	case prod.Price <= 0:
		return prod, &APIError{Code: 400, Message: "Product price must be greater than zero"}
	}

	//3. Stamp, and save to DB
	prod.CreateTs = time.Now().Unix()
	err = svc.repo.Create(&prod)
	if err != nil {
		return prod, &APIError{Code: 500, Message: fmt.Sprintf("Error creating product: %s", err.Error())}
	}

	//4. Return the product
	return prod, nil
}
