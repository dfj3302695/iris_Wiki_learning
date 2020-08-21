package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	// 另一方式 app.Use(before) 然后 app.Done(after)
	app.Get("/", before, mainHandler, after)
	app.Listen(":8080")
}

func before(ctx iris.Context) {
	shareInformation := "this is a shareable information between handles"
	// 请求路径
	requestPath := ctx.Path()
	println("before the mainHandler请求路径：" + requestPath)
	ctx.Values().Set("info", shareInformation)
	ctx.Next() // 执行下一下Handler, 这里是mainHandler
}

func after(ctx iris.Context) {
	println("after the mainHandler")
}

func mainHandler(ctx iris.Context) {
	println("Inside mainHandler")
	info := ctx.Values().GetString("info")
	ctx.HTML("<h1>Response</h1>")
	ctx.HTML("<br/>Info:" + info)
	ctx.Next() // 执行 "after" handler
}
