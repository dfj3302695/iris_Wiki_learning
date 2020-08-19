# Routing path parameter types
Iris拥有你所见过的最简单、最强大的路由过程。

Iris有自己的interpeter，用于路由的路径语法、解析和评估（是的，就像一种编程语言！）。

它的速度很快，怎么快？它计算它的需求，如果不需要任何特殊的regexp，那么它只需要用低级路径语法来注册路由，否则它就会预先编译regexp并添加必要的中间件。这意味着与其他路由器或网络框架相比，你的性能成本为零。

## Parameters
路径参数的名称应该只包含英文字母，不允许使用数字或符号，如'_'。

不要混淆ctx.Params()和ctx.Values()。

路径参数的值可以从ctx.Params()中获取。
Context的本地存储可以用于处理程序和中间件之间的通信，可以存储到ctx.Values()。
内置的可用参数类型可以在下表中找到。

| Param Type | Go Type | Validation | Retrieve Helper |
| :----- | ---- | :---- | :---- |
| :string       | string  | anything \(single path segment\)                                                                                                                                 | Params\(\)\.Get       |
| :int          | int     | \-9223372036854775808 to 9223372036854775807 \(x64\) or \-2147483648 to 2147483647 \(x32\), depends on the host arch                                             | Params\(\)\.GetInt    |
| :int8         | int8    | \-128 to 127                                                                                                                                                     | Params\(\)\.GetInt8   |
| :int16        | int16   | \-32768 to 32767                                                                                                                                                 | Params\(\)\.GetInt16  |
| :int32        | int32   | \-2147483648 to 2147483647                                                                                                                                       | Params\(\)\.GetInt32  |
| :int64        | int64   | \-9223372036854775808 to 9223372036854775807                                                                                                                     | Params\(\)\.GetInt64  |
| :uint         | uint    | 0 to 18446744073709551615 \(x64\) or 0 to 4294967295 \(x32\), depends on the host arch                                                                           | Params\(\)\.GetUint   |
| :uint8        | uint8   | 0 to 255                                                                                                                                                         | Params\(\)\.GetUint8  |
| :uint16       | uint16  | 0 to 65535                                                                                                                                                       | Params\(\)\.GetUint16 |
| :uint32       | uint32  | 0 to 4294967295                                                                                                                                                  | Params\(\)\.GetUint32 |
| :uint64       | uint64  | 0 to 18446744073709551615                                                                                                                                        | Params\(\)\.GetUint64 |
| :bool         | bool    | "1" or "t" or "T" or "TRUE" or "true" or "True" or "0" or "f" or "F" or "FALSE" or "false" or "False"                                                            | Params\(\)\.GetBool   |
| :alphabetical | string  | lowercase or uppercase letters                                                                                                                                   | Params\(\)\.Get       |
| :file         | string  | lowercase or uppercase letters, numbers, underscore \(\_\), dash \(\-\), point \(\.\) and no spaces or other special characters that are not valid for filenames | Params\(\)\.Get       |
| :path         | string  | anything, can be separated by slashes \(path segments\) but should be the last part of the route path                                                            | Params\(\)\.Get       |

#### 用法：
```go
app.Get("/users/{id:uint64}", func(ctx iris.Context){
    id := ctx.Params().GetUint64Default("id", 0)
    // [...]
})
```

| Built\-in Func   | Param Types        |
|:----------       |:------------|
| 
regexp\(expr string\)                                                                                                           | :string                                                                                               |
| 
prefix\(prefix string\)                                                                                                         | :string                                                                                               |
| 
suffix\(suffix string\)                                                                                                         | :string                                                                                               |
| 
contains\(s string\)                                                                                                            | :string                                                                                               |
| 
min\(minValue int or int8 or int16 or int32 or int64 or uint8 or uint16 or uint32 or uint64  or float32 or float64\)            | :string\(char length\), :int, :int8, :int16, :int32, :int64, :uint, :uint8, :uint16, :uint32, :uint64 |
| max\(maxValue int or int8 or int16 or int32 or int64 or uint8 or uint16 or uint32 or uint64 or float32 or float64\)              | :string\(char length\), :int, :int8, :int16, :int32, :int64, :uint, :uint8, :uint16, :uint32, :uint64 |
| 
range\(minValue, maxValue int or int8 or int16 or int32 or int64 or uint8 or uint16 or uint32 or uint64 or float32 or float64\) | :int, :int8, :int16, :int32, :int64, :uint, :uint8, :uint16, :uint32, :uint64                         |

#### 用法：
```go
app.Get("/profile/{name:alphabetical max(255)}", func(ctx iris.Context){
    name := ctx.Params().Get("name")
    // len(name) <=255 otherwise this route will fire 404 Not Found
    // and this handler will not be executed at all.
})
```
#### 自己动手:
RegisterFunc可以接受任何返回func(paramValue string)bool的函数。
或者只是一个func(string) bool。如果验证失败，那么它将发射404或任何其他关键字的状态码。
```go
latLonExpr := "^-?[0-9]{1,3}(?:\\.[0-9]{1,10})?$"
latLonRegex, _ := regexp.Compile(latLonExpr)

// Register your custom argument-less macro function to the :string param type.
// MatchString is a type of func(string) bool, so we use it as it is.
app.Macros().Get("string").RegisterFunc("coordinate", latLonRegex.MatchString)

app.Get("/coordinates/{lat:string coordinate()}/{lon:string coordinate()}",
func(ctx iris.Context) {
    ctx.Writef("Lat: %s | Lon: %s", ctx.Params().Get("lat"), ctx.Params().Get("lon"))
})
```
注册您的自定义宏函数，该函数接受两个int参数。
```go
app.Macros().Get("string").RegisterFunc("range",
func(minLength, maxLength int) func(string) bool {
    return func(paramValue string) bool {
        return len(paramValue) >= minLength && len(paramValue) <= maxLength
    }
})

app.Get("/limitchar/{name:string range(1,200) else 400}", func(ctx iris.Context) {
    name := ctx.Params().Get("name")
    ctx.Writef(`Hello %s | the name should be between 1 and 200 characters length
    otherwise this handler will not be executed`, name)
})
```
注册你的自定义宏函数，它可以接受字符串[...,...]的slice。
```go
app.Macros().Get("string").RegisterFunc("has",
func(validNames []string) func(string) bool {
    return func(paramValue string) bool {
        for _, validName := range validNames {
            if validName == paramValue {
                return true
            }
        }

        return false
    }
})

app.Get("/static_validation/{name:string has([kataras,maropoulos])}",
func(ctx iris.Context) {
    name := ctx.Params().Get("name")
    ctx.Writef(`Hello %s | the name should be "kataras" or "maropoulos"
    otherwise this handler will not be executed`, name)
})

```
#### 示例代码：
```go
func main() {
    app := iris.Default()

    // This handler will match /user/john but will not match neither /user/ or /user.
    app.Get("/user/{name}", func(ctx iris.Context) {
        name := ctx.Params().Get("name")
        ctx.Writef("Hello %s", name)
    })

    // This handler will match /users/42
    // but will not match /users/-1 because uint should be bigger than zero
    // neither /users or /users/.
    app.Get("/users/{id:uint64}", func(ctx iris.Context) {
        id := ctx.Params().GetUint64Default("id", 0)
        ctx.Writef("User with ID: %d", id)
    })

    // However, this one will match /user/john/send and also /user/john/everything/else/here
    // but will not match /user/john neither /user/john/.
    app.Post("/user/{name:string}/{action:path}", func(ctx iris.Context) {
        name := ctx.Params().Get("name")
        action := ctx.Params().Get("action")
        message := name + " is " + action
        ctx.WriteString(message)
    })

    app.Listen(":8080")
}
```
> 当参数类型缺失时，则默认为字符串类型，因此{name:string}和{name}指的是同样的东西


