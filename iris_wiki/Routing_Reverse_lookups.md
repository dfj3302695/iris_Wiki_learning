## 反向路由
正如在Routing一章中提到的，Iris提供了几种处理程序注册方法，每种方法都会返回一个Route实例。
###路由命名
路由的命名很简单，因为我们只需调用返回的*Route与Name字段来定义一个名称。  
[sample4](sample4/main.go)

###路由反转，也就是从route name生成URL。
当我们为一个特定的路径注册处理程序时，我们就可以根据我们传递给Iris的结构化数据来创建URL。在上面的例子中，
我们已经命名了三个路由器，其中一个甚至可以接受参数。如果我们使用的是默认的html/template视图引擎，
我们可以使用一个简单的操作来反转路由（并生成实际的URL）。
```go
Home: {{ urlpath "home" }}
About: {{ urlpath "about" }}
Page 17: {{ urlpath "page" "17" }}
```
[sample5](sample5/main.go)

###在代码中使用路由命名
我们可以使用以下方法/函数来处理命名的路由（及其参数）。
* GetRoutes function to get all registered routes
* GetRoute(routeName string) method to retrieve a route by name
* URL(routeName string, paramValues ...interface{}) method to generate url string based on supplied parameters
* Path(routeName string, paramValues ...interface{} method to generate just the path (without host and protocol) portion of the URL based on provided values