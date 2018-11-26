package main

import (
	_ "beego-test/routers"

	"github.com/astaxie/beego"
)

func init() {
	beego.SetLogger("file", `{
		"filename": "logs/application.log",
		"level": 7,
		"daily": true,
		"maxdays": 365,
		"rotate": true
	}`)
}

func main() {
	beego.Run()
}
