## 中间件
当我们谈论Iris中的中间件时，我们谈论的是在HTTP请求生命周期中的主处理程序代码之前和/或之后运行代码。
例如，日志中间件可能会将传入的请求细节写到日志中，然后调用处理程序代码，然后再将响应的细节写到日志中。
关于中间件的一个很酷的事情是，这些单元非常灵活和可重用。

一个中间件只是func(ctx iris.Context)的一个Handler形式，
中间件正在执行的时候，前一个中间件调用ctx.Next()，这可以用于认证，
即：如果请求认证了，那么调用ctx.Next()与请求中的其余处理程序链一起处理，否则发射一个错误响应
### 写一个中间件
代码：[sample6](sample6/main.go)

### 全局中间件
代码：[sample7](sample7/main.go)

你也可以使用ExecutionRules来强制执行Done处理程序，而不需要在你的路由处理程序中使用
ctx.Next()，像这样做。

```go
app.SetExecutionRules(iris.ExecutionRules{
    // Begin: ...
    // Main:  ...
    Done: iris.ExecutionOptions{Force: true},
})
```
示例代码[sample8](sample8/main.go)

### 转换 http.Handler/HandlerFunc
然而你并不局限于它们，你可以自由使用任何与net/http包兼容的第三方中间件。

Iris与其他软件不同，它是100%兼容的Go的net/http，这也是为什么大多数大公司将它应用到工作流程中的，
比如非常著名的US Television Network，都信任Iris的原因；它总是与最新的，标准的net/http包保持一致，
而标准的net/http包是由Go的作者在维护和更新。

任何为net/http编写的第三方中间件都可以使用`iris.FromStd(aThirdPartyMiddleware)`与Iris兼容。
记住，`ctx.ResponseWriter()`和`ctx.Request()`
返回的同样的[http.Handler](https://golang.org/pkg/net/http/#Handler) `net/http` input参数。
* [From func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)](sample9)
* From http.Handler or http.HandlerFunc
* From func(http.HandlerFunc) http.HandlerFunc
***
下面是一些专门为Iris制作的 handlers 列表。  
#### Built-in

| Middleware           | Example                                               |
|----------------------|-------------------------------------------------------|
| basic authentication | iris/_examples/auth/basicauth                         |
| request logger       | iris/_examples/logging/request-logger                 |
| HTTP method override | iris/middleware/methodoverride/methodoverride_test.go |
| profiling (pprof)    | iris/_examples/pprof                                  |
| Google reCAPTCHA     | iris/_examples/auth/recaptcha                         |
| hCaptcha             | iris/_examples/auth/recaptcha                         |
| recovery             | iris/_examples/recover                                |
| rate                 | iris/_examples/request-ratelimit                      |
| jwt                  | iris/_examples/auth/jwt                               |
| requestid            | iris/middleware/requestid/requestid_test.go           |


#### Community

| Middleware | Description | Example |
|---|---|---|
| jwt | Middleware checks for a JWT on the Authorization header on incoming requests and decodes it | iris-contrib/middleware/jwt/_example |
| cors | HTTP Access Control | iris-contrib/middleware/cors/_example |
| secure | Middleware that implements a few quick security wins | iris-contrib/middleware/secure/_example |
| tollbooth | Generic middleware to rate-limit HTTP requests | iris-contrib/middleware/tollboothic/_examples/limit-handler |
| cloudwatch | AWS cloudwatch metrics middleware | iris-contrib/middleware/cloudwatch/_example |
| new relic | Official New Relic Go Agent
 | iris-contrib/middleware/newrelic/_example |
| prometheus | Easily create metrics endpoint for the prometheus instrumentation tool | iris-contrib/middleware/prometheus/_example |
| casbin | An authorization library that supports access control models like ACL, RBAC, ABAC | iris-contrib/middleware/casbin/_examples |
| raven | Sentry client in Go | iris-contrib/middleware/raven/_example |
| csrf | Cross-Site Request Forgery Protection | iris-contrib/middleware/csrf/_example |
| go-i18n | i18n Iris Loader for nicksnyder/go-i18n | iris-contrib/middleware/go-i18n/_example |
| throttler | Rate limiting access to HTTP endpoints | iris-contrib/middleware/throttler/_example |

### 第三方中间件
Iris有自己的中间件形式func(ctx iris.Context)，但它也兼容所有net/http中间件形式。请看[这里](https://github.com/kataras/iris/tree/master/_examples/convert-handlers)。
