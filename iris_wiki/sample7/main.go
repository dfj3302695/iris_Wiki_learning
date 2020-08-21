package main

import (
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	// 注册路由
	app.Get("/", indexHandler)
	app.Get("/contact", contactHandler)

	// 顺序不重要 ，`UseGlobal` 和 `DoneGlobal` 存在的路由和将来的路由都适用
	// 记住："Use"和 "Done "适用于当前party's及其children。
	// 所以如果我们在注册路由之前使用`app.Use/Done。
	// 在这种情况下，它将像UseGlobal/DoneGlobal一样工作。
	// 因为`app`是根 "Party"。
	app.UseGlobal(before)
	app.DoneGlobal(after)

	app.Listen(":8080")
}

func before(ctx iris.Context) {
	println("before")
	ctx.Next() // wiki 少了这行 ， 不然不会到下一步indexHandler
}

func after(ctx iris.Context) {
	println("after")
}

func indexHandler(ctx iris.Context) {
	// write something to the client as a response.
	ctx.HTML("<h1>Index</h1>")
	println("index")
	ctx.Next() // execute the "after" handler registered via `Done`.
}

func contactHandler(ctx iris.Context) {
	// write something to the client as a response.
	ctx.HTML("<h1>Contact</h1>")
	println("contact")
	ctx.Next() // execute the "after" handler registered via `Done`.
}
