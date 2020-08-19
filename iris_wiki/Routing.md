# Routing
## Handle 类型
Handler，顾名思义就是处理请求  

一个Handler响应一个HTTP请求。它将回复头和数据写入Context.ResponseWriter()，然后返回。返回标志着请求已完成；
在处理程序调用完成后或与之同时使用Context是无效的。

根据HTTP客户端软件、HTTP协议版本以及客户端和iris服务器之间的任何中间件，
可能无法在写入context.ResponseWriter()之后再从Context.Request().Body中读取。
谨慎地Handles应该先读取Context.Request().Body，然后再回复

除了读取 body 外，Handles 不应修改提供的Context。

如果Handler恐慌(panics)，服务器（Handler的调用者）会认为恐慌的影响是孤立于活动请求的。它恢复恐慌(panics)，
将跟踪logs到服务器错误日志，并挂断连接。
```go
type Handler func(iris.Context)
```
一旦Handler被注册，我们可以使用返回的Route实例给处理程序注册一个名称，以便于调试，
或者在views中匹配相对路径。更多信息，
请查看[Reverse lookups](https://github.com/kataras/iris/wiki/Routing-reverse-lookups) 部分。

##行为
Iris的默认行为是接受和注册路径为/api/user的路由，没有尾部的斜杠，
如果客户端试图到达$your_host/api/user/，那么Iris路由器会自动将其重定向到$your_host/api/user，
以便由Iris路由器处理。如果客户端试图访问$your_host/api/user/，
那么Iris路由器会自动将其重定向到$your_host/api/user，以便由注册的路由来处理。这就是现代设计API的方式。

然而，如果你想禁用请求资源的路径校正，你可以将iris.WithoutPathCorrection选项的iris配置传递给你的app.Run。
例子:
```go
// [app := iris.New...]
// [...]

app.Listen(":8080", iris.WithoutPathCorrection)
```
如果你想在/api/user和/api/user/路径上保持相同的处理方法和路由，而不需要重定向(常见的情况)，
只需要使用iris.WithoutPathCorrectionRedirection选项。
```go
app.Listen(":8080", iris.WithoutPathCorrectionRedirection)
```

##API
支持所有的HTTP方法，开发者也可以在同一路径上用不同的方法注册处理程序。

第一个参数是HTTP方法，第二个参数是路由的请求路径，第三个变量参数应该包含一个或多个iris.Handler，当客户端向服务器请求该特定的resouce路径时，由注册的命令执行。

Example code:
```go
app := iris.New()

app.Handle("GET", "/contact", func(ctx iris.Context) {
    ctx.HTML("<h1> Hello from /contact </h1>")
})
```
为了方便终端开发者，iris为所有的HTTP方法提供了方法帮助器（method helpers），第一个参数是路由的请求路径，
第二个可变参数应该包含一个或多个iris.Handler，当用户从服务器请求特定的resouce路径时，由注册的命令执行。
第一个参数是路由的请求路径，第二个变量参数应该包含一个或多个iris.Handler，
当用户向服务器请求该特定的resouce路径时，由注册的命令执行。

Example code:
```go
app := iris.New()

// Method: "GET"
app.Get("/", handler)

// Method: "POST"
app.Post("/", handler)

// Method: "PUT"
app.Put("/", handler)

// Method: "DELETE"
app.Delete("/", handler)

// Method: "OPTIONS"
app.Options("/", handler)

// Method: "TRACE"
app.Trace("/", handler)

// Method: "CONNECT"
app.Connect("/", handler)

// Method: "HEAD"
app.Head("/", handler)

// Method: "PATCH"
app.Patch("/", handler)

// register the route for all HTTP Methods
app.Any("/", handler)

func handler(ctx iris.Context){
    ctx.Writef("Hello from method: %s and path: %s\n", ctx.Method(), ctx.Path())
}
```

###Offline Routes
在Iris中，有一个特殊的方法，你也可以使用。它叫做None，你可以用它来隐藏一个路由，不让外人看到，
但仍然能够通过Context.Exec方法从其他路由的处理程序中调用它。每个API Handle方法都会返回Route值。
一个Route的IsOnline方法，会回报该路由的当前状态。你可以通过它的Route.Method字段的值，
将路由的状态从离线变为在线，反之亦然。当然每次在服务时改变路由器都需要调用app.RefreshRouter()，
这样使用起来才安全。下面来看一个比较完整的例子。  
[Offline Routes 例子](./sample3/main.go)
####How to run
1. go run main.go
2. Open a browser at http://localhost:8080/invisible/iris and you'll see that you get a 404 not found error,
3. however the http://localhost:8080/execute will be able to execute that route.
4. Now, if you navigate to the http://localhost:8080/change and refresh the 
/invisible/iris tab you'll see that you can see it.

##Routes组
一组被路径前缀的路由可以（可选择）共享相同的中间件处理程序和模板布局。一个组也可以有一个嵌套组。

.Party被用来对路由进行分组，开发者可以声明无限数量的（嵌套）组。
Example code:
```go
app := iris.New()

users := app.Party("/users", myAuthMiddlewareHandler)

// http://localhost:8080/users/42/profile
users.Get("/{id:uint64}/profile", userProfileHandler)
// http://localhost:8080/users/messages/1
users.Get("/messages/{id:uint64}", userMessageHandler)
```
同样也可以用PartyFunc方法来写，它接受子路由器(Party)。
```go
app := iris.New()

app.PartyFunc("/users", func(users iris.Party) {
    users.Use(myAuthMiddlewareHandler)

    // http://localhost:8080/users/42/profile
    users.Get("/{id:uint64}/profile", userProfileHandler)
    // http://localhost:8080/users/messages/1
    users.Get("/messages/{id:uint64}", userMessageHandler)
})
```

##Path 参数
与你见过的其他路由器不同，Iris的那款路由器可以处理各种路由路径，而不会在它们之间发生冲突。  

Matches only GET "/".
```go
app.Get("/", indexHandler)
```
匹配所有以"/assets/**/*"为前缀的GET请求，它是一个通配符，
ctx.Params().Get("asset")等于/assets/之后的任何路径。
```go
app.Get("/assets/{asset:path}", assetsWildcardHandler)
```
匹配所有以"/profile/"为前缀且后面有单一路径部分的GET请求。

```go
app.Get("/profile/{username:string}", userHandler)
```
只匹配GET"/profile/me"，而且它与/profile/{username:string}或任何根号通配符/{root:path}不冲突。

```go
app.Get("/profile/me", userHandler)
```

匹配所有以/users/为前缀的GET请求，后面的数字应该等于或大于1。
```go
app.Get("/user/{userid:int min(1)}", getUserHandler)
```
匹配所有以/users/为前缀的DELETE请求，后面的数字应该等于或大于1。
```go
app.Delete("/user/{userid:int min(1)}", deleteUserHandler)
```
匹配所有GET请求，除了那些已经被其他路由处理的请求。例如，在本例中，
上述路由; /, /assets/{asset:path}, /profile/{username}, "/profile/me", /user/{userid:int ...}。
它与其他的路由并不冲突(!)。
```go
app.Get("{root:path}", rootWildcardHandler)
```
匹配所有的GET请求。

1. /u/abcd 映射到:alphabetical（如果:alphabetical已注册，否则:string）。
1. /u/42 映射到:uint (如果:uint已注册，否则:int)
2. /u/-1 映射到 :int (如果 :int 已注册，否则 :string)
1 ./u/abcd123 映射到:string。
```go
app.Get("/u/{username:string}", func(ctx iris.Context) {
	ctx.Writef("username (string): %s", ctx.Params().Get("username"))
})

app.Get("/u/{id:int}", func(ctx iris.Context) {
	ctx.Writef("id (int): %d", ctx.Params().GetIntDefault("id", 0))
})

app.Get("/u/{uid:uint}", func(ctx iris.Context) {
	ctx.Writef("uid (uint): %d", ctx.Params().GetUintDefault("uid", 0))
})

app.Get("/u/{firstname:alphabetical}", func(ctx iris.Context) {
	ctx.Writef("firstname (alphabetical): %s", ctx.Params().Get("firstname"))
})
```
匹配所有对/abctenchars.xml和/abcdtenchars的GET请求。
```go
app.Get("/{alias:string regexp(^[a-z0-9]{1,10}\\.xml$)}", PanoXML)
app.Get("/{alias:string regexp(^[a-z0-9]{1,10}$)}", Tour)
```
你可能会好奇{id:uint64}或:path或min(1)是什么。它们是(类型化的)动态路径参数，
可以在它们上面注册函数。请阅读[Path Parameter Types](Routiing%20path%20types.md)。
