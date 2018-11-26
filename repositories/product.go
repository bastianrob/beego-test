package repositories

import (
	"beego-test/models"

	"github.com/astaxie/beego/orm"
)

//IProductRepository interface to access product data
type IProductRepository interface {
	//List of products in page (starts at 1), and rpp (rows per page)
	List(page int, rpp int) ([]models.Product, error)

	//Find a product by its ID
	Find(id int64) (models.Product, error)

	//Create a new product
	Create(*models.Product) error
}

//private implementation of ProductRepository
type productRepository struct {
	db orm.Ormer
}

//NewProductRepository factory pattern to instantiate product repository
func NewProductRepository(db orm.Ormer) IProductRepository {
	return &productRepository{db: db}
}

func (repo *productRepository) List(page int, rpp int) (products []models.Product, err error) {
	limit, offset := rpp, (page-1)*rpp
	_, err = repo.db.
		Raw(`SELECT * FROM product LIMIT ? OFFSET ?`, limit, offset).
		QueryRows(&products)
	return products, err
}

func (repo *productRepository) Find(id int64) (product models.Product, err error) {
	err = repo.db.
		Raw(`SELECT * FROM product WHERE id = ? LIMIT 1`, id).
		QueryRow(&product)
	return product, err
}

func (repo *productRepository) Create(product *models.Product) error {
	_, err := repo.db.Insert(&product)
	return err
}
