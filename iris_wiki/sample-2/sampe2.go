package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()
	// 从"./views"加载所有模板
	// .html为扩展名
	// 使用常规包`html/template`

	// Configuration 的其中一种使用方法
	//app.Configure(iris.WithoutStartupLog)

	app.RegisterView(iris.HTML("./view",".html"))

	// GET 方法
	// 路径： "/"
	app.Get("/", func(ctx iris.Context){
		// 邦定  hello world 到 {{ .message }}
		ctx.ViewData("message", "hello world")
		// Render 模板文件 ./views/hello.html
		ctx.View("hello.html")
	})

	// GET 方法
	// 资源： http://localhost:8080/user/42
	//
	// 是否可以用正则代替？
	// 很容易
	// 只要参数类型是'string',可以做用`regexp`宏
	// 如：app.Get("/user/{id:string regexp(^[0-9]+$)}")
	app.Get("/user/{id:uint64}", func(ctx iris.Context) {
		userID, _ := ctx.Params().GetUint64("id")
		ctx.Writef("User ID: %d",userID)
	})

	// 使用正则
	app.Get("/username/{name:string regexp(^[0-9]+$)}", func(ctx iris.Context) {
		userName := ctx.Params().GetString("name")
		ctx.Writef("User Name:%v", userName)
	})

	// 指定port 开始服务器
	app.Listen(":8080")
}
