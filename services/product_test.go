package services_test

import (
	"beego-test/models"
	"beego-test/services"
	"errors"
	"testing"
	"time"

	"github.com/astaxie/beego/orm"

	"github.com/stretchr/testify/assert"
)

type productRepoMock struct {
	listProd []models.Product
	listErr  error

	findProd models.Product
	findErr  error

	createErr error
}

func (repo *productRepoMock) List(page int, rpp int) ([]models.Product, error) {
	return repo.listProd, repo.listErr
}

func (repo *productRepoMock) Find(id int64) (models.Product, error) {
	return repo.findProd, repo.findErr
}

func (repo *productRepoMock) Create(*models.Product) error {
	return repo.createErr
}

func Test_IProductService_List(t *testing.T) {
	assert := assert.New(t)
	mockRepo := productRepoMock{}

	t.Log("1. Expect success with 1 product listed")
	mockRepo.listProd = append([]models.Product{}, models.Product{
		ID:       1,
		CreateTs: time.Now().Unix(),
		Name:     "Product 1",
		Category: "CAT-1",
		Price:    50000,
	})
	serv := services.NewProductService(&mockRepo)
	prod, err := serv.List(1, 10)
	assert.Nil(err, "No error")
	assert.Equal(1, len(prod), "List returned is of size: 1")
	assert.Equal(int64(1), prod[0].ID, "Product ID is 1")

	//---
	t.Log("2. Expect error: page must be >= 1")
	prod, err = serv.List(0, 10)
	assert.Error(err, "Page must be greater than zero")

	//---
	t.Log("3. Expect error: rows per page must be >= 10")
	prod, err = serv.List(1, 1)
	assert.Error(err, "RPP must be greater than 10")

	//---
	t.Log("4. Expect error: something wrong with the DB, error code 500")
	mockRepo.listErr = errors.New("Something wrong with the DB")
	prod, err = serv.List(1, 10)
	assert.Error(err, "Something wrong with the DB")
	assert.Equal(500, err.Code)
}

func Test_IProductService_Find(t *testing.T) {
	assert := assert.New(t)
	mockRepo := productRepoMock{}

	t.Log("1. Expect success find product")
	mockRepo.findProd = models.Product{
		ID:       1,
		CreateTs: time.Now().Unix(),
		Name:     "Product 1",
		Category: "CAT-1",
		Price:    50000,
	}
	serv := services.NewProductService(&mockRepo)
	prod, err := serv.Find(1)
	assert.Nil(err, "No error")
	assert.Equal(int64(1), prod.ID, "Product ID is 1")

	//---
	t.Log("2. Expect error: id must be >= 1")
	prod, err = serv.Find(0) //assume db serial starts with 1
	assert.Error(err, "ID parameter must exists")

	//---
	t.Log("3. Expect error: not found")
	mockRepo.findErr = orm.ErrNoRows
	prod, err = serv.Find(1)
	mockRepo.findErr = nil
	assert.Error(err, "Error product not found")
	assert.Equal(404, err.Code)

	//---
	t.Log("4. Expect error: something wrong with the DB")
	mockRepo.findErr = errors.New("Something wrong with the DB")
	prod, err = serv.Find(1)
	mockRepo.findErr = nil
	assert.Error(err, "Error product not found")
	assert.Equal(500, err.Code)
}

func Test_IProductService_Create(t *testing.T) {
	assert := assert.New(t)
	mockRepo := productRepoMock{}

	t.Log("1. Expect success create product")
	serv := services.NewProductService(&mockRepo)
	jsonReq := `{
		"name": "Product 1",
		"category": "CAT-1",
		"price": 50000
	}`
	prod, err := serv.Create([]byte(jsonReq))
	assert.Nil(err, "No error")
	assert.Equal("Product 1", prod.Name, "Product name is Product 1")
	assert.NotEqual(0, prod.CreateTs, "Product created timestamp is filled")

	//---
	t.Log("2. Expect error: price not defined")
	jsonReq = `{
		"name": "Product 1",
		"category": "CAT-1"
	}`
	prod, err = serv.Create([]byte(jsonReq))
	assert.Error(err, "Product price must be > 0")
	assert.Equal(400, err.Code)

	//---
	t.Log("3. Expect error: category not defined")
	jsonReq = `{
		"name": "Product 1",
		"price": 50000
	}`
	prod, err = serv.Create([]byte(jsonReq))
	assert.Error(err, "Product category must be filled")
	assert.Equal(400, err.Code)

	//---
	t.Log("4. Expect error: name not defined")
	jsonReq = `{
		"category": "CAT-1",
		"price": 50000
	}`
	prod, err = serv.Create([]byte(jsonReq))
	assert.Error(err, "Product name must be filled")
	assert.Equal(400, err.Code)

	//---
	t.Log("5. Expect error: invalid JSON")
	jsonReq = `i am not a json`
	prod, err = serv.Create([]byte(jsonReq))
	assert.Error(err, "Product name must be filled")
	assert.Equal(400, err.Code)
	assert.Equal("Invalid JSON request", err.Message)

	//---
	t.Log("6. Expect error: Something wrong with the DB")
	jsonReq = `{
		"name": "Product 1",
		"category": "CAT-1",
		"price": 50000
	}`
	mockRepo.createErr = errors.New("Something wrong with the DB")
	prod, err = serv.Create([]byte(jsonReq))
	assert.Error(err, "Something wrong with the DB")
	assert.Equal(500, err.Code)
}
