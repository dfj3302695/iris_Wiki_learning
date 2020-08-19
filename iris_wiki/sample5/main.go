package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.RegisterView(iris.HTML("./view",".html"))
	// 定义一个function
	h := func(ctx iris.Context) {
		ctx.View("home.html")
	}
	// handler 注册和命名
	home := app.Get("/",h)
	home.Name = "home"

	// 或者
	app.Get("/about",h).Name = "about"
	app.Get("/page/{id}",h).Name = "page"
	app.Listen(":8080")
}
