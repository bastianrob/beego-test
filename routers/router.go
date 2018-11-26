package routers

import (
	"beego-test/controllers"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/1.0",
		beego.NSNamespace("/product",
			beego.NSRouter("/list", &controllers.ProductController{}, "get:List"),
			beego.NSRouter("/find", &controllers.ProductController{}, "get:Find"),
			beego.NSRouter("/create", &controllers.ProductController{}, "post:Create"),
		),
	)
	beego.AddNamespace(ns)
}
