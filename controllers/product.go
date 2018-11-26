package controllers

import (
	"beego-test/repositories"
	"beego-test/services"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

//ProductController entry point for product resource
type ProductController struct {
	beego.Controller
}

//List of product
//@param page	defaults at 1
//@param rpp	defaults at 20
func (c *ProductController) List() {
	page, err := c.GetInt("page", 1)
	if err != nil {
		c.CustomAbort(400, "invalid page")
	}

	rpp, err := c.GetInt("rpp", 20)
	if err != nil {
		c.CustomAbort(400, "invalid rows per page")
	}

	repo := repositories.NewProductRepository(orm.NewOrm())
	prod, apiErr := services.NewProductService(repo).List(page, rpp)
	if apiErr != nil {
		c.CustomAbort(apiErr.Code, apiErr.Message)
	}

	c.Data["json"] = prod
	c.ServeJSON()
}

//Find a product
//@param id	required
func (c *ProductController) Find() {
	id, err := c.GetInt64("id", 0)
	if err != nil {
		c.CustomAbort(400, "invalid parameter id")
	}

	repo := repositories.NewProductRepository(orm.NewOrm())
	prod, apiErr := services.NewProductService(repo).Find(id)
	if apiErr != nil {
		c.CustomAbort(apiErr.Code, apiErr.Message)
	}

	c.Data["json"] = prod
	c.ServeJSON()
}

//Create a new product
//@body	models.Product
func (c *ProductController) Create() {
	repo := repositories.NewProductRepository(orm.NewOrm())
	prod, apiErr := services.NewProductService(repo).Create(c.Ctx.Input.RequestBody)
	if apiErr != nil {
		c.CustomAbort(apiErr.Code, apiErr.Message)
	}

	c.Data["json"] = prod
	c.ServeJSON()
}
