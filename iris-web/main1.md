### 导入`iris`包

```go
import "github.com/kataras/iris"
```

`iris` 框架为您的下一个网站，`API`，或分布式应用程序提供一个简单且优美，并且易于使用的基础。

该项目的源代码和其他扩展可以从`GITHUB`获得：

> [https://github.com/kataras/iris](https://iris-go.com/donate)

##### 安装

唯一的环境要求是安装了`GO`语言环境，最低版本要求是`1.8`版本，但是强烈推荐使用`1.102`。

```go
 go get -u github.com/kataras/iris
```

示例:  `main1.go`

```go
package main

import "github.com/kataras/iris"

//绑定 数据字段 结构体 例子
type User struct {
	Username  string `json:"username"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	City      string `json:"city"`
	Age       int    `json:"age"`
}

func main() {
	app := iris.New()
	//设置数据驱动的模板引擎为 std html/template （go 语言标准模板 参见 html/template 标准库）
	// 当然 iris 还有很多数据驱动的模板引擎 以后会介绍到
	//被解析与加载的所有模板文件放在  ./web/views 文件夹里面，并且以 .html 为文件扩张
	// Reload 方法设置为 true 表示开启开发者模式 将会每一次请求都重新加载 views 文件下的所有模板
	// RegisterView 注册加载 模板文件 与加载配置
	app.RegisterView(iris.HTML("./web/views", ".html").Reload(true))
	// 为特定HTTP错误注册自定义处理程序方法
	// 当出现 StatusInternalServerError 500错误，将执行第二参数回调方法
	app.OnErrorCode(iris.StatusInternalServerError, func(ctx iris.Context) {
		// ctx.Values() 是一个很有用的东西，主要用来使 处理方法与中间件 通信 记住真的很重要
		// ctx.Values().GetString("error") 获取自定义错误提示信息
		errMessage := ctx.Values().GetString("error")
		if errMessage != "" {
			ctx.Writef("Internal server error: %s", errMessage)
			return
		}
		ctx.Writef("(Unexpected) internal server error")
	})
	// context.Handler 类型 每一个请求都会先执行此方法 app.Use(context.Handler)
	app.Use(func(ctx iris.Context) {
		ctx.Application().Logger().Infof("Begin request for path: %s", ctx.Path())
		ctx.Next()
	})
	// context.Handler 类型 每一个请求最后执行 app.Done(context.Handler)
	app.Done(func(ctx iris.Context) {})
	// Method POST: http://localhost:8080/decode
	app.Post("/decode", func(ctx iris.Context) {
		var user User
		// 请求参数格式化  请求参数是json类型转化成 User类型
		// 比如 post 参数 {username:'xxxx'} 转成 User 类型
		//把 json 类型请求参数 转成结构体
		ctx.ReadJSON(&user)
		ctx.Writef("%s %s is %d years old and comes from %s", user.Firstname, user.Lastname, user.Age, user.City)
	})
	//Method GET: http://localhost:8080/encode
	app.Get("/encode", func(ctx iris.Context) {
		doe := User{
			Username:  "Johndoe",
			Firstname: "John",
			Lastname:  "Doe",
			City:      "Neither FBI knows!!!",
			Age:       25,
		}
		//把结构体类型  转成json
		ctx.JSON(doe)
	})
	//Method GET: http://localhost:8080/profile/anytypeofstring
	app.Get("/profile/{username:string}", profileByUsername)
	//上面的 username 只接收字符串也可以用正则表达式代替并缩小字符串范围 前提要对正则有所了解
	// 类如: app.Get("/profile/{username:string regexp(^[a-zA-Z ]+$)}")

	//当不明确定义传值类型的时候，默认为字符串类型
	// app.Get("/profile/{username}", profileByUsername)等同于下面的
	//app.Get("/profile/{username:string}", profileByUsername)

	//app.Party 定义路由组  第一个参数 设置路由相同的前缀 第二个参数为中间件
	usersRoutes := app.Party("/users", logThisMiddleware)
	//两个{}只是把相同路由组的放在一个区块，没有其他用特殊含义
	{
		// Method GET: http://localhost:8080/users/42
		// 表示id必须是int类型 最小值为 1
		usersRoutes.Get("/{id:int min(1)}", getUserByID)
		// Method POST: http://localhost:8080/users/create
		usersRoutes.Post("/create", createUser)
	}
	//监听 HTTP/1.x & HTTP/2 客户端在  localhost 端口号8080 设置字符集
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8"))
}

func logThisMiddleware(ctx iris.Context) {
	//ctx.Path() 请求的url
	ctx.Application().Logger().Infof("Path: %s | IP: %s", ctx.Path(), ctx.RemoteAddr())

	//ctx.Next 继续往下一个处理方法 中间件 如果没有他 那么就不会执行 usersRoutes
	ctx.Next()
}

func profileByUsername(ctx iris.Context) {
	//获取路由参数
	username := ctx.Params().Get("username")
	//向数据模板传值 当然也可以绑定其他值
	ctx.ViewData("Username", username)
	//渲染模板 ./web/views/profile.html
	//把获得的动态数据username 绑定在 ./web/views/profile.html 模板 语法{{}} {{ .Username }}
	ctx.View("profile.html")
}

func getUserByID(ctx iris.Context) {
	//下面的可以转换成  .Values().GetInt/GetInt64
	//ctx.Values().GetInt("id")
	userID := ctx.Params().Get("id")
	user := User{Username: "username" + userID}
	//xml  输出
	ctx.XML(user)
}

func createUser(ctx iris.Context) {
	var user User
	//ctx.ReadForm 格式请求数据 与ctx.ReadJSON相似 不过接收的是 Form请求
	//记住 post 字段取名  Username 结构体字段体
	err := ctx.ReadForm(&user)
	if err != nil {
		ctx.Values().Set("error", "creating user, read and parse form failed. "+err.Error())
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}
	//{{ . }} 表示 User 的对象  取User的字段 i.e  {{ .Username }} , {{ .Firstname}} etc..
	//向数据模板传值 当然也可以绑定其他值
	ctx.ViewData("", user)
	//渲染模板 ./web/views/create_verification.html
	ctx.View("create_verification.html")
}
```

#### 监听与优雅的关闭

您可以启动服务器来侦听任何类型的 `net.Listener`或`http:server`实例。服务器的初始化方法应该在最后通过`Run`函数传递。

示例：

```go
//监听`tcp` 的 `0.0.0.0:8080` 网络地址 也就是 `net.Listener` 类型
app.Run(iris.Addr(":8080"))

//与上面之前一样的，但是使用自定义`HTTP.Server`也可能在其他地方使用达到相同的效果。

app.Run(iris.Server(&http.Server{Addr:":8080"}))

// 使用自定义的 `net.Listener`
l, err := net.Listen("tcp4", ":8080")
if err != nil {
    panic(err)
}
app.Run(iris.Listener(l))

// 使用文件 安全传输层协议`TLS`
app.Run(iris.TLS("127.0.0.1:443", "mycert.cert", "mykey.key"))

// 自动 安全传输层协议`TLS`
app.Run(iris.AutoTLS(":443", "example.com", "admin@example.com"))

// UNIX 套接字
if errOs := os.Remove(socketFile); errOs != nil && !os.IsNotExist(errOs) {
    app.Logger().Fatal(errOs)
}

l, err := net.Listen("unix", socketFile)

if err != nil {
    app.Logger().Fatal(err)
}

if err = os.Chmod(socketFile, mode); err != nil {
    app.Logger().Fatal(err)
}

app.Run(iris.Listener(l))

//使用任何 func（）error，
//启动监听者的责任取决于你这个方式，
//为了简单起见，我们将使用
//`net/http`包的ListenAndServe函数。
app.Run(iris.Raw(&http.Server{Addr:":8080"}).ListenAndServe)
```
UNIX和BSD主机可以利用重用端口功能

示例:  `main2.go`

```go
package main

import (
	//tcplisten包 提供各种可自定义的TCP net.Listener 与 性能相关的选项 Linux 特性
	//第一 SO_REUSEPORT。 此选项允许线性扩展服务器性能 在多CPU服务器上。
	//关详细信息，请参阅 https://www.nginx.com/blog/socket-sharding-nginx-release-1-9-1/
	//第二 TCP_DEFER_ACCEPT。 此选项期望服务器从接受的读取写入之前的连接
	//第三 TCP_FASTOPEN 关详细信息，请参阅https://lwn.net/Articles/508865/。

	"github.com/valyala/tcplisten"
	"github.com/kataras/iris"
)
// 安装 tcplisten $ go get github.com/valyala/tcplisten
// $ go run main.go

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello World!</b>")
	})
	//对用上面的三个选项
	listenerCfg := tcplisten.Config{
		ReusePort:   true,
		DeferAccept: true,
		FastOpen:    true,
	}
	l, err := listenerCfg.NewListener("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	app.Run(iris.Listener(l))
}
```
这就是监听，你可以在需要时完全控制。

让我们继续学习如何捕获`CONTROL + C / COMMAND + C`或`unix kill`命令并优雅地关闭服务器。

```go
正确关闭使用CONTROL + C / COMMAND + C或者当发送的 kill命令是默认可用时
```

为了手动管理应用程序中断时要执行的操作，我们必须使用选项 `WithoutInterruptHandler`禁用默认行为并注册新的中断处理程序（全局，跨所有可能的主机）。

说白了就是自定义响应 `CONTROL + C / COMMAND + C` 的函数，也就是自定义关闭服务

示例：  `main3.go`

```go
package main

import (
	stdContext "context"
	"time"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		// 关闭所有的 host
		app.Shutdown(ctx)
	})

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(" <h1>hi, I just exist in order to see if the server is closed</h1>")
	})

	// http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithoutInterruptHandler)
}
```

#### 主机 Hosts

在`Run`方法之后，可以通过`Application＃Hosts`字段访问为您的应用程序提供服务的所有主机。

但最常见的情况是您可能需要在`Run`方法之前访问主机，有两种方式可以访问主机主管，如下所示：

第一种方法是使用`app.NewHost`创建一个新主机，并使用其`Serve`或`Listen`函数之一通过`iris #Raw` Runner启动应用程序。 请注意，这种方式需要额外导入`net / http`包。

示例：
```go
h := app.NewHost(&http.Server{Addr:":8080"})
h.RegisterOnShutdown(func(){
    println("terminate")
})

app.Run(iris.Raw(h.ListenAndServe))
```
其次，更简单的方法是使用`host.Configurator`。

请注意，当使用go <1.9时，此方法需要额外的导入语句 `github.com/kataras/iris/core/host`，如果您的目标是go1.9，那么您可以使用`iris #Supervisor`并省略 额外的主机导入。

我们之前看到的所有常见的 **Runners**（`iris＃Addr，iris＃Listener，iris #Server，iris＃TLS，iris #AutoTLS`）接受`host.Configurator`的可变参数，只有`func（* host）.Supervisor）`。 因此，`Application`使您有权通过这些修改自动创建的主机主管。

示例：  `main4.go`
```go
package main

import (
	stdContext "context"
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
	"fmt"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML("<h1>Hello, try to refresh the page after ~10 secs</h1>")
	})

	app.Logger().Info("Wait 10 seconds and check your terminal again")
	//在这里模拟一个关机动作......
	go func() {
		<-time.After(10 * time.Second)
		timeout := 5 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		//关闭所有主机，这将通知我们已注册的回调 在configureHost func中
		app.Shutdown(ctx)
	}()
	//像往常一样启动服务器，唯一的区别就是
	//我们正在添加第二个（可选）功能
	//配置刚刚创建的主机管理。
	// http://localhost:8080
	//等待10秒钟并检查您的终端
	app.Run(iris.Addr(":8080", configureHost), iris.WithoutServerError(iris.ErrServerClosed))
}

func configureHost(su *host.Supervisor) {
	//这里我们可以完全访问将要创建的主机 下面三个方法 
	//在`Run`函数中 我们注册了一个关闭“事件”回调

	// RegisterOnShutdown注册一个函数来调用Shutdown。
	//这可以用于正常关闭具有的连接
	//经过NPN / ALPN协议升级或被劫持。
	//此函数应启动特定于协议的正常关闭，
	//但不应该等待关机完成。
	su.RegisterOnShutdown(func() {
		println("terminate")
	})
	//su.RegisterOnError 当出现错误则调用此方法
	//su.RegisterOnServe 启动服务 调用此方法
	 su.RegisterOnServe(func(taskHost host.TaskHost) {
		fmt.Print( taskHost.HostURL())
		println("test")
	 })
}
```
通过下面链接，阅读有关监听和正常关机的更多信息：

> [https://github.com/kataras/iris/tree/master/_examples/#http-listening](https://github.com/kataras/iris/tree/master/_examples/#http-listening)

#### 路由

支持所有HTTP方法，开发人员还可以为不同方法注册相同路径的处理程序。 

* 第一个参数是`HTTP`方法  `GET` `HEAD` `POST` `DELETE` `PUT` `CONNECT` `OPTIONS` `TRACE`
* 第二个参数是请求路径，也就url,路由的请求路径
* 第三个可变参数应该包含一个或多个`iris.Handler` 响应路径的具体操作

当用户从服务器请求该特定的资源路径时，由注册的方法执行。

示例：
```go
app := iris.New()

app.Handle("GET", "/contact", func(ctx iris.Context) {
    ctx.HTML("<h1> Hello from /contact </h1>")
})
```
为了使用户更容易，`iris`为所有`HTTP`方法提供了相应的函数。
* 第一个参数是路由的请求路径
* 第二个可变参数应该包含一个或多个iris.Handler

当用户从服务器请求该特定的资源路径时，由注册的方法执行。

示例：
```go
//这个例子完美的解释了，上面的定义

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

// 注册一个所有 `http` 方法 都可以调用的路由
app.Any("/", handler)

func handler(ctx iris.Context){
    ctx.Writef("Hello from method: %s and path: %s", ctx.Method(), ctx.Path())
}
```
#### 分组路由

按路径前缀分组的一组路由可以（可选）共享相同的中间件处理程序和模板布局。 一个组也可以有一个嵌套组。

`.Party`用于分组路由，开发人员可以声明无限数量的（嵌套）组。

示例：
```go
users := app.Party("/users", myAuthMiddlewareHandler)

// http://myhost.com/users/42/profile
users.Get("/{id:int}/profile", userProfileHandler)

// http://myhost.com/users/messages/1
users.Get("/messages/{id:int}", userMessageHandler)
```
#### 自定义HTTP错误

`iris`开发人员能够注册自己的处理程序，如`404`未找到的`http`状态，`500`内部服务器错误等等。

示例：
```go
//当报 `404` 时候渲染自定义的 404 错误模板
app.OnErrorCode(iris.StatusNotFound, func(ctx iris.Context){
    ctx.View("errors/404.html")
})

app.OnErrorCode(500, func(ctx iris.Context){
    // 编写对500 错误的响应
})
```

#### 基本HTTP API

在`iris`的表现主义路由器的帮助下，您可以安全地构建任何形式的API。

示例： `main5.go`
```go
package main
import "github.com/kataras/iris"

func main() {
	app := iris.New()

	//状态代码注册404（错误）自定义处理程序，
	//未找到路由时触发，或者通过ctx.StatusCode（iris.StatusNotFound）手动触发。
	app.OnErrorCode(iris.StatusNotFound, notFoundHandler)

	// GET -> HTTP 方法
	// /   -> 路径
	// func(ctx iris.Context) -> 响应路由处理程序

	//第三个接收器（参数）应包含路由的处理程序，它们按顺序执行。
	app.Handle("GET", "/", func(ctx iris.Context) {
		//阅读  GOPATH/src/github.com/kataras/iris/context/context.go
		//概述所有上下文的方法（有很多这些方法，阅读它们，你将学习iris如何工作）
		ctx.HTML("Hello from " + ctx.Path()) // Hello from /
	})

	app.Get("/home", func(ctx iris.Context) {
		ctx.Writef(`Same as app.Handle("GET", "/", [...])`)
	})

	app.Get("/donate", donateHandler, donateFinishHandler)

	//不要忘记更多“匹配路由”的动态路径示例！
	app.Get("/api/users/{userid:int min(1)}", func(ctx iris.Context) {
		userID, err := ctx.Params().GetInt("userid")

		if err != nil {
			ctx.Writef("error while trying to parse userid parameter," +
				"this will never happen if :int is being used because if it's not integer it will fire Not Found automatically.")
			ctx.StatusCode(iris.StatusBadRequest)
			return
		}

		ctx.JSON(map[string]interface{}{
			//当然，您可以传递任何自定义返回值。
			"user_id": userID,
		})
	})
	//下面内容  上面已经解释过论
	// app.Post("/", func(ctx iris.Context){}) -> for POST http method.
	// app.Put("/", func(ctx iris.Context){})-> for "PUT" http method.
	// app.Delete("/", func(ctx iris.Context){})-> for "DELETE" http method.
	// app.Options("/", func(ctx iris.Context){})-> for "OPTIONS" http method.
	// app.Trace("/", func(ctx iris.Context){})-> for "TRACE" http method.
	// app.Head("/", func(ctx iris.Context){})-> for "HEAD" http method.
	// app.Connect("/", func(ctx iris.Context){})-> for "CONNECT" http method.
	// app.Patch("/", func(ctx iris.Context){})-> for "PATCH" http method.
	// app.Any("/", func(ctx iris.Context){}) for all http methods.


	//多个路由可以包含具有不同http映射方法的相同路径。
	// 您可以捕获任何路径创建错误：
	// route, err := app.Get(...)
	// 为路由设置名称 route: route.Name = "myroute"

	//您还可以按路径前缀，共享中间件和完成处理程序对路由进行分组。
	adminRoutes := app.Party("/admin", adminMiddleware)

	adminRoutes.Done(func(ctx iris.Context) {
		//如果ctx.Next（），则始终执行
		ctx.Application().Logger().Infof("response sent to " + ctx.Path())
	})
	// adminRoutes.Layout("/views/layouts/admin.html") /设置这些路由的视图布局，请参阅视图示例中的更多内容。

	// GET: http://localhost:8080/admin
	adminRoutes.Get("/", func(ctx iris.Context) {
		// [...]
		ctx.StatusCode(iris.StatusOK) // 默认 200 == iris.StatusOK
		ctx.HTML("<h1>Hello from admin/</h1>")

		ctx.Next()//为了执行adminRoutes.Done() 处理程序
	})

	// GET:  http://localhost:8080/admin/login
	adminRoutes.Get("/login", func(ctx iris.Context) {
		// [...] 对用操作
	})
	// POST: http://localhost:8080/admin/login
	adminRoutes.Post("/login", func(ctx iris.Context) {
		// [...] 对用操作
	})

	// subdomains, easier than ever, should add localhost or 127.0.0.1 into your hosts file,
	// etc/hosts on unix or .
	//子域名，比以往更容易，应该将localhost或127.0.0.1添加到您的hosts文件中，
	// unix 上的 /etc/hosts 或 Windows上的 C:/windows/system32/drivers/etc/hosts on windows。
	v1 := app.Party("v1.")
	{ //大括号是可选的，它只是样式的类型，可以直观地分组路由。

		// http://v1.localhost:8080
		v1.Get("/", func(ctx iris.Context) {
			ctx.HTML("Version 1 API. go to <a href='" + ctx.Path() + "/api" + "'>/api/users</a>")
		})

		//可以无限嵌套 .Party
		usersAPI := v1.Party("/api/users")
		{
			// http://v1.localhost:8080/api/users
			usersAPI.Get("/", func(ctx iris.Context) {
				ctx.Writef("All users")
			})
			// http://v1.localhost:8080/api/users/42
			usersAPI.Get("/{userid:int}", func(ctx iris.Context) {
				ctx.Writef("user with id: %s", ctx.Params().Get("userid"))
			})
		}
	}

	//通配符子域。
	wildcardSubdomain := app.Party("*.")
	{
		wildcardSubdomain.Get("/", func(ctx iris.Context) {
			ctx.Writef("Subdomain can be anything, now you're here from: %s", ctx.Subdomain())
		})
	}

	// http://localhost:8080
	// http://localhost:8080/home
	// http://localhost:8080/donate
	// http://localhost:8080/api/users/42
	// http://localhost:8080/admin
	// http://localhost:8080/admin/login
	//
	// http://localhost:8080/api/users/0
	// http://localhost:8080/api/users/blabla
	// http://localhost:8080/wontfound
	//
	// if hosts edited:
	//  http://v1.localhost:8080
	//  http://v1.localhost:8080/api/users
	//  http://v1.localhost:8080/api/users/42
	//  http://anything.localhost:8080
	app.Run(iris.Addr(":8080"))
}

func adminMiddleware(ctx iris.Context) {
	// [...]
	ctx.Next()  //移动到下一个处理程序，如果你有任何auth逻辑，则不要这样做
}

func donateHandler(ctx iris.Context) {
	ctx.Writef("Just like an inline handler, but it can be " +
		"used by other package, anywhere in your project.")

	//ctx.Values().Set()
	//让我们将值传递给下一个处理程序
	//值是处理程序（或中间件）在彼此之间进行通信的方式。
	ctx.Values().Set("donate_url", "https://github.com/kataras/iris#-people")
	ctx.Next() //为了执行链中的下一个处理程序(donateFinishHandler)。
}

func donateFinishHandler(ctx iris.Context) {
	//值可以是任何类型的对象，因此我们可以将值转换为字符串
	//但是iris提供了一个很容易做到的，如果没有定义donate_url，那么它会返回一个空字符串。
	donateURL := ctx.Values().GetString("donate_url")
	ctx.Application().Logger().Infof("donate_url value was: " + donateURL)
	ctx.Writef("\n\nDonate sent(?).")
}

func notFoundHandler(ctx iris.Context) {
	ctx.HTML("Custom route for 404 not found http code, here you can render a view, html, json <b>any valid response</b>.")
}
```
#### 参数化路径

在前面的示例中，我们已经看到了静态路由，路由组，子域，通配符子域，带有单个已知参数和自定义http错误的参数化路径的小示例，现在是时候查看通配符参数和宏。

iris，就像 net/http std包由Handler注册路由的处理程序一样，iris的处理程序类型只是一个func（ctx iris.Context），其中的上下文来自github.com/kataras/iris/context。

iris 拥有您遇到过的最简单，最强大的路由过程。

同时，iris 有自己的路径（就像编程语言一样）用于路径的路径语法及其动态路径参数解析和评估，我们称它们为“宏”用于快捷方式。 怎么样？ 
它计算了它的需求，如果不需要任何特殊的正则表达式，那么它只是用低级路径语法注册路由，否则它预先编译正则表达式并添加必要的中间件。

参数的标准宏类型：
```go
+------------------------+
| {param:string}         |
+------------------------+
string type   //字符串类型
anything      //任何参数

+------------------------+
| {param:int}            |
+------------------------+
int type               //int 类型
only numbers (0-9)     //仅0-9 组合

+------------------------+
| {param:long}           |
+------------------------+
int64 type              //int int64
only numbers (0-9)      //仅0-9 组合

+------------------------+
| {param:boolean}        |
+------------------------+
bool type               //布尔类型

only "1" or "t" or "T" or "TRUE" or "true" or "True"
or "0" or "f" or "F" or "FALSE" or "false" or "False"

+------------------------+
| {param:alphabetical}   |
+------------------------+
alphabetical/letter type            // 字母类型
letters only (upper or lowercase)  //大写或者小写

+------------------------+
| {param:file}           |
+------------------------+
file type                      //文件类型 可传下面类型
letters (upper or lowercase)  //大写或者小写
numbers (0-9)                 //数字类型
underscore (_)               //下划线
dash (-)                     //分割线
point (.)                     // 点
no spaces ! or other character 不能是 空格 ！ 和其他字符串

+------------------------+
| {param:path}           |
+------------------------+
path type     //路径类型
例如: /path1/path2/path3 , ctx.Params().Get("param") == "/path1/path2/path3"
```

如果缺少类型，则参数的类型默认为字符串，因此{param} == {param：string}。

如果在该类型上找不到函数，则使用“字符串”的类型函数。即：

```go
{param:int min(3)}
```
除了`iris`提供基本类型和一些默认的“宏功能”这一事实，你也可以自己注册！

注册命名路径参数功能：
```go
app.Macros().Int.RegisterFunc("min", func(argument int) func(paramValue string) bool {
    [...]
    //返回 true 或者 false  true 表示验证通过
    return true/false -> true means valid.
})
```

在`func（argument...）`你可以有任何标准类型，它将在服务器启动之前进行验证，所以不关心这里的性能，它在服务时运行的唯一一件事就是返回`func（paramValue string）bool`。

```go
{param:string equal(iris)} , //下面将完成
app.Macros().String.RegisterFunc("equal", func(argument string) func(paramValue string) bool {
    return func(paramValue string){ return argument == paramValue }
})
```
示例：  `main6.go`
```go
package main

import (
	"github.com/kataras/iris"
	"strconv"
)

func main()  {
	app := iris.New()

	//您可以使用“string”类型，该类型对于可以是任何内容的单个路径参数有效。
	app.Get("/username/{name}", func(ctx iris.Context) {
		ctx.Writef("Hello %s", ctx.Params().Get("name"))
	}) //缺少类型= {name：string}

	//让我们注册我们的第一个附加到int宏类型的宏。
	// "min" = the function 方法
	// "minValue" = 是 该方法的参数
	//func（string）bool = 宏的路径参数赋值器，这在服务时执行
	//用户使用min（...）宏参数函数请求包含：int宏类型的路径。
	app.Macros().Int.RegisterFunc("min", func(minValue int) func(string) bool {
		//在此之前做任何事情[...]

		//在这种情况下，我们不需要做任何事情
		return func(paramValue string) bool {
			n, err := strconv.Atoi(paramValue)
			if err != nil {
				return false
			}
			return n >= minValue
		}
	})
	// http://localhost:8080/profile/id>=1
	//这将抛出404，如果请求路由：/ profile/0，/profile/blabla，/profile/-1
	//宏参数函数当然是可选的。
	app.Get("/profile/{id:int min(1)}", func(ctx iris.Context) {
		//第二个参数是错误，但它总是为 nil，因为我们使用宏，
		//验证已经发生。
		id, _ := ctx.Params().GetInt("id")
		ctx.Writef("Hello id: %d", id)
	})
	//更改每个路由的宏响应的错误代码程序：
	//验证错误 则由默认的抛出404 改成 504
	app.Get("/profile/{id:int min(1)}/friends/{friendid:int min(1) else 504}", func(ctx iris.Context) {
		id, _ := ctx.Params().GetInt("id")
		friendid, _ := ctx.Params().GetInt("friendid")
		ctx.Writef("Hello id: %d looking for friend id: ", id, friendid)
	}) //如果没有传递所有路由的宏，这将抛出e 504错误代码而不是404。

	// http://localhost:8080/game/a-zA-Z/level/0-9
	//http://localhost:8080/game/aA/level/12
	// 记住，alphabetical 只是小写或大写字母。
	app.Get("/game/{name:alphabetical}/level/{level:int}", func(ctx iris.Context) {
		ctx.Writef("name: %s | level: %s", ctx.Params().Get("name"), ctx.Params().Get("level"))
	})

	//让我们使用一个简单的自定义regexp来验证单个路径参数
	//它的值只是小写字母。
	// http://localhost:8080/lowercase/anylowercase
	app.Get("/lowercase/{name:string regexp(^[a-z]+)}", func(ctx iris.Context) {
		ctx.Writef("name should be only lowercase, otherwise this handler will never executed: %s", ctx.Params().Get("name"))
	})

	// http://localhost:8080/single_file/app.js
	app.Get("/single_file/{myfile:file}", func(ctx iris.Context) {
		ctx.Writef("file type validates if the parameter value has a form of a file name, got: %s", ctx.Params().Get("myfile"))
	})

	// http://localhost:8080/myfiles/any/directory/here/
	//这是唯一接受任意数量路径段的宏类型。
	app.Get("/myfiles/{directory:path}", func(ctx iris.Context) {
		//path type accepts any number of path segments, path after /myfiles/ is: any/directory/here
		ctx.Writef("path type accepts any number of path segments, path after /myfiles/ is: %s", ctx.Params().Get("directory"))
	})
	app.Run(iris.Addr(":8080"))
}
```

路径参数名称应仅包含字母，符号，包含“_”，并且不允许使用数字。 如果路由未能注册，如果你没有在.Handle / .Get上捕获第二个返回值（错误），应用程序将会没有任何警告而发错误.... 上面说过的捕获错误

最后，不要将ctx.Values（）与ctx.Params（）混淆。 路径参数的值转到ctx.Params（）和上下文的本地存储，可用于在处理程序和中间件之间进行通信转到ctx.Values（），路径参数和其余的任何自定义值都是为了你自己的功能实现而分开的。

#### **run** 执行代码

```go
$ go run main.go
```

#### 静态文件加载

```go
// StaticServe函数 可以将目录作为Web资源提供 路径助手函数
//它是 Static*函数的最简单形式    Static*  代表以 Static开头的函数，下面将详细介绍
//与StaticWeb的用法几乎相同
//只接受一个必需参数，即systemPath，
//将使用相同的路径来注册GET和HEAD方法路由。
//如果第二个参数为空，否则requestPath是第二个参数
//它使用gzip压缩（每个请求压缩，没有文件缓存）。

// Returns the GET *Route.  返回GET * Route
StaticServe(systemPath string, requestPath ...string) (*Route, error)

// StaticContent将GET和HEAD方法路由注册到requestPath
//准备好提供原始静态字节，缓存内存。
//
// Returns the GET *Route.  返回GET * Route
StaticContent(reqPath string, cType string, content []byte) (*Route, error)

//当文件在app可执行文件中分发时使用的StaticEmbedded，主要使用go-bindata
//第一个参数是请求路径，即vdir中的文件将被提供给的路径，例如“/ static”
//第二个参数是（虚拟）目录路径，例如“./assets”
//第三个参数是Asset函数
//第四个参数是AssetNames函数。
//注： go-bindata 是一个文件管理函数集合  当你参考了下面的github地址就更容易明白了

// Returns the GET *Route.   返回GET * Route
//
// Example: https://github.com/kataras/iris/tree/master/_examples/file-server/embedding-files-into-app
StaticEmbedded(requestPath string, vdir string, assetFn func(name string) ([]byte, error), namesFn func() []string) (*Route, error)

// Favicon提供静态图标
//接受2个参数，第二个是可选的
// favPath（string），声明 __.ico 的系统目录路径
// requestPath（string），它是路由的路径，默认情况下这是“/favicon.ico”，因为有些浏览器首先会尝试默认获取，
//如果您有多个favicon（桌面，移动设备等），您可以声明自己的路径

//这个函数会为你添加一个路由，它会将/yuorpath/yourfile.ico静态地提供给/yourfile.ico
//（没有什么特别的，你自己无法处理）。
//请注意，您必须在必须自动提供的每个favicon（桌面，移动设备等）上调用它。

// Returns the GET *Route.  返回GET * Route
Favicon(favPath string, requestPath ...string) (*Route, error)

// StaticWeb returns a handler that serves HTTP requests
// with the contents of the file system rooted at directory.
// StaticWeb返回一个服务于HTTP请求的处理程序
//使用以目录为根的文件系统的内容。

//第一个参数：路由路径
//第二个参数：系统文件目录
//第三个OPTIONAL参数：异常路由
//（=优先考虑这些路由而不是静态处理程序）
//查看更多选项 app.StaticHandler。

//app.StaticWeb("/static", "./static")  
//上面第二参数 表示静态文件在 main.go 同级目录 static 里面
//第一个参数表示请求路由为static  例如 请求路由是 host/static/xxx.js

//作为一种特殊情况，返回的文件服务器会重定向任何请求
//以“/index.html”结尾到同一路径，没有最终版本
//“index.html”

// StaticWeb调用StaticHandler（systemPath，listingDirectories：false，gzip：false）。
//
// Returns the GET *Route.  返回GET * Route
StaticWeb(requestPath string, systemPath string, exceptRoutes ...*Route) (*Route, error)
```
示例：  `main7.go`
```go
package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()
	//现在可以直接访问  localhost:8080/favicon.ico
	app.Favicon("./web/ico/one.ico")

	// app.Favicon("./web/ico/one.ico", "/favicon.ico")
	// 上面可以 这样访问  localhost:8080/favicon.ico

	app.Get("/", func(ctx iris.Context) {
		ctx.HTML(`<a href="/favicon.ico"> press here to see the favicon.ico</a>.
        At some browsers like chrome, it should be visible at the top-left side of the browser's window,
        because some browsers make requests to the /favicon.ico automatically,
        so iris serves your favicon in that path too (you can change it).`)
	}) //如果您没有显示favicon，请尝试清除浏览器的缓存.

	app.Run(iris.Addr(":8080"))
```
更多示例可以在这里找到：[https://github.com/kataras/iris/tree/master/_examples/beginner/file-server](https://github.com/kataras/iris/tree/master/_examples/beginner/file-server)

#### 中间件生态系统

中间件只是有序处理程序链的概念。 中间件可以在全局，每组，每子域和每个路由中注册。

示例：

```go
// 全局中间件
//在任何路由之前，都会先执行此中间件
app.Use(func(ctx iris.Context){
   // ... 编写你在请求每一个路由前想做的操作 是否登录判断 token 验证

   ctx.Next() //为了继续下一个处理程序，
   //如果缺少那么，链处理程序中的下一个将不会被执行，
   //对身份验证中间件很有用
})

// 全局中间件
//在任何路由之前或之后，将中间件预先安排到所有路由
app.UseGlobal(handler1, handler2, handler3)

// 当前路由 中间件  针对路由
app.Post("/login", authenticationHandler, loginPageHandler)

// 当前组  针对每组
users := app.Party("/users", usersMiddleware)
users.Get("/", usersIndex)

//当前子域名  针对子域名
mysubdomain := app.Party("mysubdomain.", firstMiddleware)
mysubdomain.Use(secondMiddleware)
mysubdomain.Get("/", mysubdomainIndex)

//每个通配符，动态子域
dynamicSub := app.Party(".*", firstMiddleware, secondMiddleware)
dynamicSub.Get("/", func(ctx iris.Context){
  ctx.Writef("Hello from subdomain: "+ ctx.Subdomain())
})
```

iris能够将您以前使用的任何外部第三方处理程序包装并转换为Web应用程序。 让我们转换[https://github.com/rs/cors](https://github.com/rs/cors) net/http外部中间件，它返回一个`next form`处理程序。

示例： `main8.go`

```go
package main

import (
	"github.com/rs/cors"
	"github.com/kataras/iris"
)

func main() {

	//跨与请求 下面代码表示 http://foo.com 站点下的 ajax 可以跨域请求 localhost:8080 接口
	//当http://foo.com 为 × 表示所有域名都可以请求

	//AllowedOrigins 该字段是必须的。
	// 它的值要么是请求时Origin字段的值，要么是一个*，表示接受任意域名的请求。

    //AllowCredentials  该字段可选。它的值是一个布尔值，表示是否允许发送Cookie。
    //默认情况下，Cookie不包括在CORS请求之中。
    //设为true，即表示服务器明确许可，Cookie可以包含在请求中，一起发给服务器。
    //这个值也只能设为true，如果服务器不要浏览器发送Cookie，删除该字段即可。

    //ExposedHeaders 该字段可选。CORS请求时，XMLHttpRequest对象的getResponseHeader()方法只能拿到6个基本字段该字段可选。
    //6个基本字段：Cache-Control、Content-Language、Content-Type、Expires、Last-Modified、Pragma。
    //如果想拿到其他字段，就必须在 ExposedHeaders 里面指定。
	app := iris.New()
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://foo.com"},
		AllowCredentials: true,
	}

	corsWrapper := cors.New(corsOptions).ServeHTTP

	app.WrapRouter(corsWrapper)

	v1 := app.Party("/api/v1")
	{
		v1.Get("/", h)
		v1.Put("/put", h)
		v1.Post("/post", h)
	}

	app.Run(iris.Addr(":8080"))
}

func h(ctx iris.Context) {
	ctx.Application().Logger().Infof(ctx.Path())
	ctx.Writef("Hello from %s", ctx.Path())
}
```

#### 视图引擎

Iris支持开箱即用的5个模板引擎，开发人员仍然可以使用任何外部golang模板引擎，因为`context/context＃ResponseWriter()`是一个`io.Writer`。

所有这五个模板引擎都具有通用API的共同特征，如布局，模板功能，特定组的布局，部分渲染等。

```go
标准的html
它的模板解析器是 golang.org/pkg/html/template/

Django
它的模板解析器是 github.com/flosch/pongo2

Pug(Jade)
它的模板解析器是 github.com/Joker/jade

Handlebars
它的模板解析器是 github.com/aymerick/raymond

Amber
它的模板解析器是 github.com/eknkc/amber
```

示例： `main9.go`

```go
package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()

	// - 标准 html      | iris.HTML(...)
	// - django         | iris.Django(...)
	// - pug(jade)      | iris.Pug(...)
	// - handlebars     | iris.Handlebars(...)
	// - amber          | iris.Amber(...)

	tmpl := iris.HTML("./templates", ".html")
	tmpl.Reload(true) //在每个请求上重新加载模板（开发模式）
	// 默认模板功能是：
	//
	// - {{ urlpath "mynamedroute" "pathParameter_ifneeded" }} 链接
	// - {{ render "header.html" }}   //引入其他页面
	// - {{ render_r "header.html" }} //当前页面的部分相对路径
	// - {{ yield }}   占位符 layout 页面
	// - {{ current }} 

	// 注册自定义模板func。
	tmpl.AddFunc("greet", func(s string) string {
		return "Greetings " + s + "!"
	})

	//将视图引擎注册到视图，这将加载模板。
	app.RegisterView(tmpl)

	app.Get("/", hi)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"), iris.WithCharset("UTF-8")) //默认为这个但您可以更改它。
}

func hi(ctx iris.Context) {
	ctx.ViewData("Title", "Hi Page")
	ctx.ViewData("Name", "iris") // {{.Name}} 将输出: iris
	// ctx.ViewData("", myCcustomStruct{})
	ctx.View("hi.html")
}
```

视图引擎也支持捆绑[https://github.com/shuLhan/go-bindata](https://github.com/shuLhan/go-bindata)模板文件。 go-bindata为您提供了两个函数，asset和assetNames，可以使用`.Binary`函数将这些函数设置为每个模板引擎。

程序需要加载模板文件生成输出的 HTML。在程序部署的时候，除了发布应用可执行文件外，还需要发布依赖的静态资源文件。这给发布过程添加了一些麻烦。既然发布单独一个可执行文件是非常简单的操作，就有人会想办法把静态资源文件打包进 Go 的程序文件中。

示例： `main10.go`

```go
package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()
	// $ go get -u github.com/shuLhan/go-bindata/...  切记安装时不要忘记三个点 ...
	// $ go-bindata ./templates/...
	// $ go build
	// $ ./embedding-templates-into-app
	// 不使用html文件，可以删除文件夹并运行示例
	app.RegisterView(iris.HTML("./web/views/main10", ".html").Binary(Asset, AssetNames))
	app.Get("/", hi)

	// http://localhost:8080
	app.Run(iris.Addr(":8080"))
}

type page struct {
	Title, Name string
}

func hi(ctx iris.Context) {
	ctx.ViewData("", page{Title: "Hi Page", Name: "iris"})
	ctx.View("hi.html")
}
```

一个真实的例子可以在这里找到：[https：//github.com/kataras/iris/tree/master/_examples/view/embedding-templates-into-app](https：//github.com/kataras/iris/tree/master/_examples/view/embedding-templates-into-app)。

在每个请求上启用模板的自动重新加载。 在开发人员处于开发模式时很有用，因为他们不需要在每个模板编辑时重新启动他们的应用程序。

示例： 

```go
pugEngine := iris.Pug("./templates", ".jade")
pugEngine.Reload(true) //<---设置为true以在每个请求上重新构建模板
app.RegisterView(pugEngine)
```

注意：

如果你想知道，视图引擎背后的代码来自`github.com/kataras/iris/view`包，对引擎变量的访问也可以通过`github.com/kataras/iris`包来获得。

```go
iris.HTML(...)      简写         view.HTML(...)
iris.Django(...)     >> >>      view.Django(...)
iris.Pug(...)        >> >>      view.Pug(...)
iris.Handlebars(...) >> >>      view.Handlebars(...)
iris.Amber(...)      >> >>      view.Amber(...)
```

这些模板引擎中的每一个都有不同的选项：[https：//github.com/kataras/iris/tree/master/view](https：//github.com/kataras/iris/tree/master/view)

#### Sessions

此示例将说明如何存储和访问会话中的数据。

您不需要任何第三方库，但如果您需要，您可以使用任何兼容的会话管理器。

在此示例中，我们将仅允许通过了身份验证的用户在 `/secret`页面上查看我们的信息。 要获得访问权限，首先必须访问`/login`以获取有效的会话`cookie`，并将其登录。此外，他可以访问 `/logout` 以撤消对我们的信息的访问权限。

示例： `main11.go`

```go
package main

import (
	"github.com/kataras/iris"

	"github.com/kataras/iris/sessions"
)

var (
	cookieNameForSessionID = "mycookiesessionnameid"
	sess                   = sessions.New(sessions.Config{Cookie: cookieNameForSessionID})
)

func secret(ctx iris.Context) {

	//检查用户是否经过身份验证
	if auth, _ := sess.Start(ctx).GetBoolean("authenticated"); !auth {
		ctx.StatusCode(iris.StatusForbidden)
		return
	}

	//打印信息
	ctx.WriteString("The cake is a lie!")
}

func login(ctx iris.Context) {
	session := sess.Start(ctx)

	//身份验证就在这里
	// ...

	//将用户设置为已验证
	session.Set("authenticated", true)
}

func logout(ctx iris.Context) {
	session := sess.Start(ctx)

	//撤消用户身份验证
	session.Set("authenticated", false)
}

func main() {
	app := iris.New()

	app.Get("/secret", secret)
	app.Get("/login", login)
	app.Get("/logout", logout)

	app.Run(iris.Addr(":8080"))
}
```

运行示例：

```go
$ go get github.com/kataras/iris/sessions
$ go run main.go

$ curl -s http://localhost:8080/secret
Forbidden

$ curl -s -I http://localhost:8080/login
Set-Cookie: mycookiesessionnameid=MTQ4NzE5Mz...

$ curl -s --cookie "mycookiesessionnameid=MTQ4NzE5Mz..." http://localhost:8080/secret
The cake is a lie!
```

可以使用一个（或多个）`sessiondb`来实现会话持久性。

示例： `main12.go`

```go
package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
)

//使用redis 3.0.503版测试。
//对于Windows，请参阅：https：//github.com/ServiceStack/redis-windows
func main() {
	//替换正在运行的redis'服务器设置：
	db := redis.New(service.Config{
		Network:     service.DefaultRedisNetwork,
		Addr:        service.DefaultRedisAddr,
		Password:    "",
		Database:    "",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      ""}) //可选择在redis服务器之间配置网桥

	//当control + C / cmd + C时关闭连接
	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	defer db.Close() //如果应用程序出错，请关闭数据库连接。

	sess := sessions.New(sessions.Config{
		Cookie:  "sessionscookieid",
		Expires: 45 * time.Minute}, // <= 0意味着不会过期。 默认为0。
	)

	// 重要：
	sess.UseDatabase(db)

	//其余代码保持不变
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx iris.Context) {
		s := sess.Start(ctx)
		//设置会话值
		s.Set("name", "iris")

		//测试获取
		ctx.Writef("All ok session value of the 'name' is: %s", s.GetString("name"))
	})

	app.Get("/set/{key}/{value}", func(ctx iris.Context) {
		key, value := ctx.Params().Get("key"), ctx.Params().Get("value")
		s := sess.Start(ctx)
		//设置会话值
		s.Set(key, value)

		//测试获取
		ctx.Writef("All ok session value of the '%s' is: %s", key, s.GetString(key))
	})

	app.Get("/get", func(ctx iris.Context) {
		//获取一个特定的键，作为字符串，如果没有找到只返回一个空字符串
		name := sess.Start(ctx).GetString("name")

		ctx.Writef("The 'name' on the /set was: %s", name)
	})

	app.Get("/get/{key}", func(ctx iris.Context) {
		//获取一个特定的键，作为字符串，如果没有找到只返回一个空字符串
		name := sess.Start(ctx).GetString(ctx.Params().Get("key"))

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx iris.Context) {
		//删除特定 key
		sess.Start(ctx).Delete("name")
	})

	app.Get("/clear", func(ctx iris.Context) {
		//删除所有条目
		sess.Start(ctx).Clear()
	})

	app.Get("/destroy", func(ctx iris.Context) {
		// 销毁，删除整个会话数据和cookie
		sess.Destroy(ctx)
	})

	app.Get("/update", func(ctx iris.Context) {
		//使用新日期更新过期日期
		sess.ShiftExpiration(ctx)
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}
```

更多例子：[https://github.com/kataras/iris/tree/master/_examples/sessions](https://github.com/kataras/iris/tree/master/_examples/sessions)

#### Websockets

在此示例中，我们将通过浏览器在Web套接字之间创建一个小聊天。

服务器代码示例： `main13.go`

```go
package main

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
)

func main() {
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.ServeFile("websockets.html", false) //第二个参数：是否启用gzip
	})

	setupWebsocket(app)
	// x2
	// http://localhost:8080
	// http://localhost:8080
	//写点东西，按提交，查看结果。
	app.Run(iris.Addr(":8080"))
}

func setupWebsocket(app *iris.Application) {
	//创建我们的echo websocket服务器
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	})
	ws.OnConnection(handleConnection)

	//在端点上注册服务器。
	//在websockets.html中查看内联javascript代码，此端点用于连接服务器。
	app.Get("/echo", ws.Handler())

	//提供javascript built'n客户端库，
	//请参阅websockets.html脚本标记，使用此路径。
	app.Any("/iris-ws.js", func(ctx iris.Context) {
		ctx.Write(websocket.ClientSource)
	})
}

func handleConnection(c websocket.Connection) {
	//从浏览器中读取事件
	c.On("chat", func(msg string) {
		//将消息打印到控制台，c .Context（）是 iris 的http上下文。
		fmt.Printf("%s sent: %s\n", c.Context().RemoteAddr(), msg)
		//将消息写回客户端消息所有者：
		// c.Emit("chat", msg)
		c.To(websocket.Broadcast).Emit("chat", msg)
	})
}
```

示例客户端（javascript）代码：

```html
<!-- websockets.html -->
<input id="input" type="text" />
<button onclick="send()">Send</button>
<pre id="output"></pre>
<script src="/iris-ws.js"></script>
<script>
    var input = document.getElementById("input");
    var output = document.getElementById("output");

    //ws.js 需要自己下载
    //ws 自动服务的'/iris-ws.js'
    var socket = new Ws("ws://localhost:8080/echo");
    socket.OnConnect(function () {
        output.innerHTML += "Status: Connected\n";
    });

    socket.OnDisconnect(function () {
        output.innerHTML += "Status: Disconnected\n";
    });

   //从服务器读取事件
    socket.On("chat", function (msg) {
        addMessage(msg)
    });

    function send() {
        addMessage("Me: "+input.value) //写自己
        socket.Emit("chat", input.value);//将聊天事件数据发送到websocket服务器
        input.value = ""; //清除输入
    }

    function addMessage(msg) {
        output.innerHTML += msg + "\n";
    }
</script>
```

运行示例：

```go
$ go get github.com/kataras/iris/websocket
$ go run main.go
$ start http://localhost:8080
```

#### MVC - 模型视图控制器

Iris拥有对MVC模式的一流支持，你不会在Go世界的任何其他地方找到这些东西。

示例： `main14.go`

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
)

func main() {
	app := iris.New()
	//（可选）添加两个内置处理程序
	//可以从任何与http相关的panics中恢复
	//并将请求记录到终端。
	app.Use(recover.New())
	app.Use(logger.New())

	//根据根路由器“/”提供控制器。
	mvc.New(app).Handle(new(ExampleController))

	// http://localhost:8080
	// http://localhost:8080/ping
	// http://localhost:8080/hello
	// http://localhost:8080/custom_path
	app.Run(iris.Addr(":8080"))
}

// ExampleController提供 “/”，“/ ping”和“/ hello”  路由
type ExampleController struct{}

// Get serves
// Method:   GET
// Resource: http://localhost:8080
func (c *ExampleController) Get() mvc.Result {
	return mvc.Response{
		ContentType: "text/html",
		Text:        "<h1>Welcome</h1>",
	}
}

// GetPing serves
// Method:   GET
// Resource: http://localhost:8080/ping
func (c *ExampleController) GetPing() string {
	return "pong"
}

// GetHello serves
// Method:   GET
// Resource: http://localhost:8080/hello
func (c *ExampleController) GetHello() interface{} {
	return map[string]string{"message": "Hello Iris!"}
}

// GetUserBy 服务 //方法：GET //资源路径：http://localhost:8080/user/{username:string}
// By 是一个保留的“关键字”告诉框架
// 绑定 函数输入参数中的路径参数，它也有助于在同一个控制器中使用“Get”和“GetBy”。
// func（c * ExampleController）GetUserBy(username string) mvc.Result { return mvc.View { Name：“user /username.html”， Data：username，}
```

可以使用多个，工厂将确保为此控制器的每个路由注册正确的http方法，如果需要，请取消注释：

```go
func (c *ExampleController) Post() {}
func (c *ExampleController) Put() {}
func (c *ExampleController) Delete() {}
func (c *ExampleController) Connect() {}
func (c *ExampleController) Head() {}
func (c *ExampleController) Patch() {}
func (c *ExampleController) Options() {}
func (c *ExampleController) Trace() {}

func (c *ExampleController) All() {}
//        OR
func (c *ExampleController) Any() {}

func (c *ExampleController) BeforeActivation(b mvc.BeforeActivation) {
     // 1 - > HTTP方法
     // 2 - >请求的路径
     // 3 - >此控制器的方法名称应该是该路由的处理程序。
    b.Handle("GET", "/mypath/{param}", "DoIt", optionalMiddlewareHere...)
}

//AfterActivation 后，所有依赖项都被设置 - 因此只读访问它们
//但仍可以添加自定义控制器或简单的标准处理程序。
func (c *ExampleController) AfterActivation(a mvc.AfterActivation) {}
```
**Iris Web**框架以最快的执行速度支持Request数据，模型，持久性数据和绑定。

**特点：**

支持所有`HTTP`方法，例如，如果想要提供`GET`，那么控制器应该有一个名为`Get()`的函数，你可以定义多个方法函数在同一个Controller中提供。

通过`BeforeActivation`自定义事件回调，每个控制器，将自定义控制器的struct的方法注册为具有自定义路径的处理程序（即使使用正则表达参数化路径）

示例： `main15.go`

```go
package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

func main() {
	app := iris.New()
	mvc.Configure(app.Party("/root"), myMVC)
	app.Run(iris.Addr(":8080"))
}

func myMVC(app *mvc.Application) {
	// app.Register(...)
	// app.Router.Use/UseGlobal/Done(...)
	app.Handle(new(MyController))
}

type MyController struct {}

func (m *MyController) BeforeActivation(b mvc.BeforeActivation) {
	// b.Dependencies().Add/Remove
	// b.Router().Use/UseGlobal/Done //以及您已经知道的任何标准API调用

	// 1-> http 请求方法
	// 2-> 请求路径
	// 3-> 控制器的函数名称将被解析为处理程序
	// 4-> 任何应该在MyCustomHandler之前运行的处理程序
	b.Handle("GET", "/something/{id:long}", "MyCustomHandler", anyMiddleware...)
}

// GET: http://localhost:8080/root
func (m *MyController) Get() string { return "Hey" }

// GET: http://localhost:8080/root/something/{id:long}
func (m *MyController) MyCustomHandler(id int64) string { return "MyCustomHandler says Hey" }
```

Controller结构中的持久性数据（在请求之间共享数据），通过定义对依赖项的服务或具有`Singleton`控制器范围。

共享控制器之间的依赖关系或在父MVC应用程序上注册它们，并能够在Controller内的`BeforeActivation`可选事件回调中修改每个控制器的依赖关系.

示例： 

```go
func(c *MyController) BeforeActivation(b mvc.BeforeActivation) { b.Dependencies().Add/Remove(...) }
```

访问`Context`作为控制器的字段（没有手动绑定是neede），即`Ctx iris.Context`或通过方法的输入参数

示例： 

```go
func(ctx iris.Context, otherArguments...)
```

Controller结构中的模型（在Method函数中设置并由View呈现）。 您可以从控制器的方法返回模型，或者在请求生命周期中设置字段，并在同一请求生命周期中将该字段返回到另一个方法。


像以前一样流，mvc应用程序有自己的`Router`，它是一种`iris/router.Party`，标准的iris api。 “控制器”可以注册到任何  `Party`，包括子域，Party`的开始和完成处理程序按预期工作。

可选的`BeginRequest（ctx）`函数在方法执行之前执行任何初始化，对调用中间件或许多方法使用相同的数据集合很有用。

可选的`EndRequest（ctx）`函数，用于在执行任何方法后执行任何终结。

会话动态依赖性通过管理器的“Start”绑定到MVC应用程序

例如： 

```go
mvcApp.Register(sessions.New(sessions.Config{Cookie: "iris_session_id"}).Start)
```

继承，递归。

Access to the dynamic path parameters via the controller's methods' input arguments, no binding is needed. When you use the Iris' default syntax to parse handlers from a controller, you need to suffix the methods with the `By` word, uppercase is a new sub path. Example:

通过控制器方法的输入参数访问动态路径参数，不需要绑定。 当您使用`Iris`的默认语法来解析来自控制器的处理程序时，您需要使用`By`字来为方法添加后缀，大写是一个新的子路径。

注册一个或多个相对路径并能够获得路径参数。

示例： 

```go
If `mvc.New(app.Party("/user")).Handle(new(user.Controller))`

- `func(*Controller) Get()` - `GET:/user` , 照常.
- `func(*Controller) Post()` - `POST:/user`, 照常.
- `func(*Controller) GetLogin()` - `GET:/user/login`
- `func(*Controller) PostLogin()` - `POST:/user/login`
- `func(*Controller) GetProfileFollowers()` - `GET:/user/profile/followers`
- `func(*Controller) PostProfileFollowers()` - `POST:/user/profile/followers`
- `func(*Controller) GetBy(id int64)` - `GET:/user/{param:long}`
- `func(*Controller) PostBy(id int64)` - `POST:/user/{param:long}`

If `mvc.New(app.Party("/profile")).Handle(new(profile.Controller))`

- `func(*Controller) GetBy(username string)` - `GET:/profile/{param:string}`

If `mvc.New(app.Party("/assets")).Handle(new(file.Controller))`

- `func(*Controller) GetByWildard(path string)` - `GET:/assets/{param:path}`

If `mvc.New(app.Party("/equality")).Handle(new(profile.Equality))`

- `func(*Controller) GetBy(is bool)` - `GET:/equality/{param:boolean}`
- `func(*Controller) GetByOtherBy(is bool, otherID int64)` - `GET:/equality/{paramfirst:boolean}/other/{paramsecond:long}`

//方法函数接收器支持的类型：int，int64，bool和string。
```

通过输出参数响应，也就是返回类型，可选如下

```go
func(c *ExampleController) Get() string |
(string, string) |
(string, int) |
(string, error) |
int |
(int, string) |
(any, int) |
error |
(int, error) |
(customStruct, error) |
(any, error) |
bool |
(any, bool)
customStruct |
(customStruct, int) |
(customStruct, string) |
`Result` or (`Result`, error)
```

`any`意味着一切，从自定义结构到标准语言的类型。 `Result`是一个只包含该函数的接口：Dispatch（ctx iris.Context）和Get HTTP Method函数（Post，Put，Delete ...）。

#### `Iris MVC`方法结果

`Iris`有一个非常强大和极快的`MVC`支持，你可以从方法函数返回任何类型的任何值，它将按预期发送到客户端

*如果`string`那么它就是返回 `body`。
*如果`string`是第二个输出参数，那么它就是`content`类型。
*如果`int`那么它是状态代码。
*如果`bool`为false，则通过跳过其他所有内容而抛出404找不到http错误。
*如果`error`而不是nil那么（任何类型）响应将被省略，并且将呈现带有400错误请求的错误文本。
*如果`（int，error）`并且错误不是nil那么响应结果将是错误的文本，状态代码为`int`。
*如果`custom struct`或`interface {}`或`slice`或`map`则它将呈现为json，除非跟随`string`内容类型。
*如果`mvc.Result`然后它执行它的`Dispatch`函数，那么可以使用好的设计模式在需要的地方拆分模型的逻辑。

可以在以下位置找到具有良好模式但不打算在生产中使用的示例：

[https://github.com/kataras/iris/tree/master/_examples/#mvc](https://github.com/kataras/iris/tree/master/_examples/#mvc)

#### 使用`Iris MVC`进行代码重用

通过创建彼此独立的组件，开发人员能够在其他应用程序中快速轻松地重用组件。 对于具有不同数据的另一个应用程序，可以为一个应用程序重构相同（或类似）的视图，因为视图只是处理数据如何显示给用户。

如果您不熟悉后端Web开发，请首先阅读有关MVC架构模式的内容，一个好的开始是维基百科文章：

[https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller)

#### 基础知识学习资料

但是你现在应该对框架有一个基本的想法，我们只是略微划清界限。 如果您喜欢刚刚看到的内容并希望了解更多内容，请点击以下链接：

[https://github.com/kataras/iris/tree/master/_examples](https://github.com/kataras/iris/tree/master/_examples)

中间件：

*[https://github.com/kataras/iris/tree/master/middleware](https://github.com/kataras/iris/tree/master/middleware)
*[https://github.com/iris-contrib/middleware](https://github.com/iris-contrib/middleware)

主页：   [https://iris-go.com](https://iris-go.com)

书籍（正在进行中）：  [https://docs.iris-go.com](https://docs.iris-go.com)

### 框架数据类型介绍

#### 常量

```go
const (
    StatusContinue           = 100 // RFC 7231, 6.2.1
    StatusSwitchingProtocols = 101 // RFC 7231, 6.2.2
    StatusProcessing         = 102 // RFC 2518, 10.1

    StatusOK                   = 200 // RFC 7231, 6.3.1
    StatusCreated              = 201 // RFC 7231, 6.3.2
    StatusAccepted             = 202 // RFC 7231, 6.3.3
    StatusNonAuthoritativeInfo = 203 // RFC 7231, 6.3.4
    StatusNoContent            = 204 // RFC 7231, 6.3.5
    StatusResetContent         = 205 // RFC 7231, 6.3.6
    StatusPartialContent       = 206 // RFC 7233, 4.1
    StatusMultiStatus          = 207 // RFC 4918, 11.1
    StatusAlreadyReported      = 208 // RFC 5842, 7.1
    StatusIMUsed               = 226 // RFC 3229, 10.4.1

    StatusMultipleChoices  = 300 // RFC 7231, 6.4.1
    StatusMovedPermanently = 301 // RFC 7231, 6.4.2
    StatusFound            = 302 // RFC 7231, 6.4.3
    StatusSeeOther         = 303 // RFC 7231, 6.4.4
    StatusNotModified      = 304 // RFC 7232, 4.1
    StatusUseProxy         = 305 // RFC 7231, 6.4.5

    StatusTemporaryRedirect = 307 // RFC 7231, 6.4.7
    StatusPermanentRedirect = 308 // RFC 7538, 3

    StatusBadRequest                   = 400 // RFC 7231, 6.5.1
    StatusUnauthorized                 = 401 // RFC 7235, 3.1
    StatusPaymentRequired              = 402 // RFC 7231, 6.5.2
    StatusForbidden                    = 403 // RFC 7231, 6.5.3
    StatusNotFound                     = 404 // RFC 7231, 6.5.4
    StatusMethodNotAllowed             = 405 // RFC 7231, 6.5.5
    StatusNotAcceptable                = 406 // RFC 7231, 6.5.6
    StatusProxyAuthRequired            = 407 // RFC 7235, 3.2
    StatusRequestTimeout               = 408 // RFC 7231, 6.5.7
    StatusConflict                     = 409 // RFC 7231, 6.5.8
    StatusGone                         = 410 // RFC 7231, 6.5.9
    StatusLengthRequired               = 411 // RFC 7231, 6.5.10
    StatusPreconditionFailed           = 412 // RFC 7232, 4.2
    StatusRequestEntityTooLarge        = 413 // RFC 7231, 6.5.11
    StatusRequestURITooLong            = 414 // RFC 7231, 6.5.12
    StatusUnsupportedMediaType         = 415 // RFC 7231, 6.5.13
    StatusRequestedRangeNotSatisfiable = 416 // RFC 7233, 4.4
    StatusExpectationFailed            = 417 // RFC 7231, 6.5.14
    StatusTeapot                       = 418 // RFC 7168, 2.3.3
    StatusUnprocessableEntity          = 422 // RFC 4918, 11.2
    StatusLocked                       = 423 // RFC 4918, 11.3
    StatusFailedDependency             = 424 // RFC 4918, 11.4
    StatusUpgradeRequired              = 426 // RFC 7231, 6.5.15
    StatusPreconditionRequired         = 428 // RFC 6585, 3
    StatusTooManyRequests              = 429 // RFC 6585, 4
    StatusRequestHeaderFieldsTooLarge  = 431 // RFC 6585, 5
    StatusUnavailableForLegalReasons   = 451 // RFC 7725, 3

    StatusInternalServerError           = 500 // RFC 7231, 6.6.1
    StatusNotImplemented                = 501 // RFC 7231, 6.6.2
    StatusBadGateway                    = 502 // RFC 7231, 6.6.3
    StatusServiceUnavailable            = 503 // RFC 7231, 6.6.4
    StatusGatewayTimeout                = 504 // RFC 7231, 6.6.5
    StatusHTTPVersionNotSupported       = 505 // RFC 7231, 6.6.6
    StatusVariantAlsoNegotiates         = 506 // RFC 2295, 8.1
    StatusInsufficientStorage           = 507 // RFC 4918, 11.5
    StatusLoopDetected                  = 508 // RFC 5842, 7.2
    StatusNotExtended                   = 510 // RFC 2774, 7
    StatusNetworkAuthenticationRequired = 511 // RFC 6585, 6
)
```

在`IANA`注册的`HTTP`状态代码。 请参阅：[http：//www.iana.org/assignments/http-status-codes/http-status-codes.xhtml](http：//www.iana.org/assignments/http-status-codes/http-status-codes.xhtml)来自`net/http` std包的原始副本，以便用户重新获取`net/http`的导入路径。

复制自`net/http`包

```go
const (
    MethodGet     = "GET"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodHead    = "HEAD"
    MethodPatch   = "PATCH"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
)
```

从`net/http`复制的HTTP方法

```go
const MethodNone = "NONE"
```

`MethodNone`是`iris`一种用于存储“离线”路由的特定的“虚拟”方法。

```go
const NoLayout = view.NoLayout
```

`NoLayout`禁用特定模板文件的布局`view＃NoLayout`的快捷方式。

#### 变量

```go
var (
    // HTML视图引擎
    // Conversion for the view.HTML.
    HTML = view.HTML
    // Django 视图引擎
    // Conversion for the view.Django.
    Django = view.Django
    // Handlebars 视图引擎
    // Conversion for the view.Handlebars.
    Handlebars = view.Handlebars
    // Pug 视图引擎
    // Conversion for the view.Pug.
    Pug = view.Pug
    // Amber 视图引擎
    // Conversion for the view.Amber.
    Amber = view.Amber
)
```

```go
var (
    // LimitRequestBodySize是一个设置请求体大小限制的中间件
    //对于链中的所有`next handlers` 处理程序。
    //
    //`context＃LimitRequestBodySize`的快捷方式。
    LimitRequestBodySize = context.LimitRequestBodySize
    // StaticEmbeddedHandler返回一个可以提供服务的Handler
    //嵌入到可执行文件中。
    //
    // 例子: https://github.com/kataras/iris/tree/master/_examples/file-server
    StaticEmbeddedHandler = router.StaticEmbeddedHandler
    // StripPrefix返回一个服务于HTTP请求的处理程序
    //从请求URL的Path中删除给定的前缀
    // 并调用处理程序 h。 StripPrefix处理一个
    //请求不以前缀为开头的路径
    // 回复HTTP 404未找到错误。
    
    //示例:
    // fileserver := Party#StaticHandler("./static_files", false, false)
    // h := iris.StripPrefix("/static", fileserver)
    // app.Get("/static/{f:path}", h)
    // app.Head("/static/{f:path}", h)
    StripPrefix = router.StripPrefix
    // Gzip是一个支持写作的中间件
    //使用gzip压缩，如果客户端支持。
    //
    //`context＃Gzip`的快捷方式。
    Gzip = context.Gzip
    //FromStd将原生http.Handler，http.HandlerFunc和func（w，r，next）转换为context.Handler。
    //
    //支持的表单类型：
    // 		 .FromStd(h http.Handler)
    // 		 .FromStd(func(w http.ResponseWriter, r *http.Request))
    // 		 .FromStd(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
    //
    //`handlerconv＃FromStd`的快捷方式。
    FromStd = handlerconv.FromStd
    // Cache是一个提供服务器端缓存功能的中间件
    //可以用作：`app.Get（“/”，iris.Cache，aboutHandler）`。
    //它应该在静态方法之后使用。
    //请参阅`iris＃Cache304`以获得更快捷的替代方法。
    //示例: https://github.com/kataras/iris/tree/master/_examples/#caching
    Cache = cache.Handler
    // NoCache是一个覆盖 header 头 Cache-Control，Pragma和Expires的中间件
    //在浏览器的后退和前进功能期间禁用缓存。
    //在HTML路由上很好地利用了这个中间件; 甚至在“后退”和“前进”浏览器的箭头按钮上刷新页面。
    //请参阅`iris＃StaticCache`以了解相反的行为。
    //
    //`cache＃NoCache`的快捷方式
    NoCache = cache.NoCache
    // StaticCache中间件，用于通过向客户端发送“Cache-Control”和“Expires”标头来缓存静态文件。
    //它接受一个输入参数，“cacheDur”，一个用于计算到期时间的time.Duration。
    
    //如果“cacheDur”<= 0，则返回“NoCache”中间件，用于禁用浏览器“后退”和“前进”操作之间的缓存。
    
    // 示例: `app.Use(iris.StaticCache(24 * time.Hour))` or `app.Use(iris.Staticcache(-1))`.
    // 中间件也是一个简单的Handler，也可以在另一个处理程序中调用, 示例:
    // cacheMiddleware := iris.StaticCache(...)
    // func(ctx iris.Context){
    //  cacheMiddleware(ctx)
    //  [...]
    // }
    //
   //`cache＃StaticCache`的快捷方式
    StaticCache = cache.StaticCache
    // Cache304每当发送一个`StatusNotModified`（304）
    //“If-Modified-Since”请求头（时间）在之前
    // time.Now() + expiresEvery (始终与其UTC值进行比较).
    //使用这个，`chache＃Cache304`的快捷方式，而不是“github.com/kataras/iris/cache”或iris.Cache
    //为了更好的表现。
    //与http RCF兼容的客户端（所有浏览器和postman等工具）
    //将处理缓存。
    //使用它而不是服务器端缓存的唯一缺点
    //这个方法会发送304状态代码而不是200
    //所以，如果你将它与其他微服务并排使用
    //您必须检查该状态代码以及有效的响应。

    //开发人员可以自由扩展此方法的行为
    //通过观察系统目录手动更改并使用`ctx.WriteWithExpiration`
    //基于文件修改日期的“modtime”，
    // simillary到`StaticWeb`（发送状态OK（200）和浏览器磁盘缓存而不是304）。
    //“cache＃Cache304”的快捷方式。
    Cache304 = cache.Cache304
    // CookiePath是一个`CookieOption`。
    //用它来更改cookie的Path字段。
    //
    //`context＃CookiePath`的快捷方式。
    CookiePath = context.CookiePath
    // CookieCleanPath是一个`CookieOption`。
    //用它来清除cookie的Path字段，与`CookiePath（“”）`完全相同。
    //
    //`context #CookieCleanPath`的快捷方式。
    CookieCleanPath = context.CookieCleanPath
    // CookieExpires是一个`CookieOption`。
    //通过传递cookie的生命周期，使用它来更改cookie的Expires和MaxAge字段。
    //
    //`context #CookieExpires`的快捷方式。
    CookieExpires = context.CookieExpires
    //使用它将cookie的HttpOnly字段设置为false或true。
    //对于`RemoveCookie`和`SetCookieKV`，HttpOnly字段默认为true。
    //`context #CookieHTTPOnly`的快捷方式
    CookieHTTPOnly = context.CookieHTTPOnly
    // CookieEncode是一个`CookieOption`。
    //在添加cookie时提供编码功能。
    //接受`context＃CookieEncoder`并将cookie的值设置为编码值。
    //用户是`context #SetCookie`和`context #SetCookieKV`。
    //示例: https://github.com/kataras/iris/tree/master/_examples/cookies/securecookie 
    //`context＃CookieEncode`的快捷方式。
    CookieEncode = context.CookieEncode    
    // CookieDecode是一个`CookieOption`。
    //在检索cookie时提供解码功能。
    //接受`context #CookieDecoder`并在返回'GetCookie`之前将cookie的值设置为解码值。
    //用户是`context＃GetCookie`。
    //示例: https://github.com/kataras/iris/tree/master/_examples/cookies/securecookie
    //`context＃CookieDecode`的快捷方式。
    CookieDecode = context.CookieDecode
)
```
```go
var ErrServerClosed = http.ErrServerClosed
```

在调用**Shutdown**或**Close**之后，服务器的`Serve`，`ServeTLS`，`ListenAndServe`和`ListenAndServeTLS`方法返回`ErrServerClosed`。

`http＃ErrServerClosed`的快捷方式:

```go
var RegisterOnInterrupt = host.RegisterOnInterrupt
```

`RegisterOnInterrupt`注册一个全局函数，在按下**CTRL + C/CMD + C**或收到**unix kill**命令时调用。

`host＃RegisterOnInterrupt`的快捷方式

```go
var (
    //Version是Iris Web Framework的当前版本号。
    Version = maintenance.Version
)
``` 

```go
var WithFireMethodNotAllowed = func(app *Application) {
    app.config.FireMethodNotAllowed = true
}
``` 

`WithFireMethodNotAllowed`使`FireMethodNotAllowed`设置变为可用。

**参见配置:**

```go
var WithGlobalConfiguration = func(app *Application) {
    app.Configure(WithConfiguration(YAML(globalConfigurationKeyword)))
}
``` 

`WithGlobalConfiguration`将从主目录加载全局yaml配置文件，它将整个应用程序的配置设置/覆盖到该文件的内容。 全局配置文件可由用户修改，并由多个`iris`实例使用。

当我们运行共享相同配置的多个iris服务器时，这非常有用，即使在“其他”字段中使用自定义值也是如此。

用法：`app.Configure（iris.WithGlobalConfiguration）`或`app.Run（[iris.Runner]，iris.WithGlobalConfiguration）`。

```go
var WithOptimizations = func(app *Application) {
    app.config.EnableOptimizations = true
}
``` 

`WithOptimizations`可以强制应用程序进行优化，以便在可能的情况下实现最佳性能。

**参见配置:**

```go
var WithPathEscape = func(app *Application) {
    app.config.EnablePathEscape = true
}
``` 

`WithPathEscape`使`PathEscape`设置变为`enanbles`。

**参见配置:**

```go
var WithoutAutoFireStatusCode = func(app *Application) {
    app.config.DisableAutoFireStatusCode = true
}
``` 

`WithoutAutoFireStatusCode`禁用`AutoFireStatusCode`设置。

**参见配置:**

```go
var WithoutBanner = WithoutStartupLog
``` 

`WithoutBanner`是`WithoutStartupLog`选项的转换。

当主服务器打开时，关闭一次发送到终端的信息。

```go
var WithoutBodyConsumptionOnUnmarshal = func(app *Application) {
    app.config.DisableBodyConsumptionOnUnmarshal = true
}
``` 

`WithoutBodyConsumptionOnUnmarshal`禁用`BodyConsumptionOnUnmarshal`设置。

**参见配置:**

```go
var WithoutInterruptHandler = func(app *Application) {
    app.config.DisableInterruptHandler = true
}
``` 

当按下`control / cmd + C`时，`WithoutInterruptHandler`禁用自动正常服务器关闭。

```go
var WithoutInterruptHandler = func(app *Application) {
    app.config.DisableInterruptHandler = true
}
``` 

`WithoutPathCorrection`禁用`PathCorrection`设置。

**参见配置:**

```go
var WithoutStartupLog = func(app *Application) {
    app.config.DisableStartupLog = true
}
``` 

当主服务器打开时，`WithoutStartupLog`关闭信息发送一次到终端。

```go
var WithoutVersionChecker = func(app *Application) {
    app.config.DisableVersionChecker = true
}
``` 

`WithoutVersionChecker`将禁用版本检查程序和更新程序。 如果将其传递给`Run`函数，`Iris`服务器将不会收到自动更新。 仅在准备好生产环境时使用它。

#### type Application

```go
type Application struct {
    //路由嵌入式| 公开APIBuilder和Router的公共API。
    *router.APIBuilder
    *router.Router
    ContextPool *context.Pool

    //主机包含运行此应用程序的所有服务器（主机主管）的列表。
    //
    //只有当应用程序使用`iris.Raw`选项运行程序运行（`app.Run`）时，主机才可能为空
    //否则它包含一个主机（`app.Hosts [0]`）。
    //
    //可以通过手动调用`app.NewHost`将其他主持人监督添加到该列表中。
    //
    //在“Run”或“NewHost”之后可以使用Hosts字段。
    Hosts []*host.Supervisor
    //包含已过滤或未导出的字段
}
``` 

应用程序负责管理应用程序的状态。 它包含并处理创建快速Web服务器所需的所有部分。

**[func Default](https://github.com/kataras/iris/blob/master/iris.go#L182)**

```go
func Default() *Application
``` 

默认返回一个新的**Application**实例，与`New`不同，它会在`panics`上恢复并记录传入的`http`请求。

**[func New](https://github.com/kataras/iris/blob/master/iris.go#L163)**

```go
func New() *Application
``` 
`New`创建并返回一个新的空`iris * Application`实例。

**[func (*Application) Build](https://github.com/kataras/iris/blob/master/iris.go#L753)**

```go
func (app *Application) Build() error
``` 

构建一次，构建框架。 它使用默认宏和非常接近`iris`的模板函数构建默认路由器。

**[func (*Application) ConfigurationReadOnly](https://github.com/kataras/iris/blob/master/iris.go#L240)**

```go
func (app *Application) ConfigurationReadOnly() context.ConfigurationReadOnly
``` 

`ConfigurationReadOnly`返回一个不允许字段写入的对象。

**[func (*Application) Configure](https://github.com/kataras/iris/blob/master/iris.go#L231)**

```go
func (app *Application) Configure(configurators ...Configurator) *Application
``` 

当需要修改框架实例时，可以调用**Configure**。 它接受框架实例并返回一个错误，如果它不是**nil**则打印到记录器。 有关更多信息，请参阅**configuration.go**。

返回自身以便像`app:= New().Configure(...)`

**[func (*Application) ConfigureHost](https://github.com/kataras/iris/blob/master/iris.go#L495)**

```go
func (app *Application) ConfigureHost(configurators ...host.Configurator) *Application
``` 

**ConfigureHost**接受一个或多个`host #Configuration`，这些配置器函数可以访问由`app.Run`创建的主机，当应用程序初始好，提供公有，它们将被执行。

这是与`app.Run`自动创建的主机交互的另一种方式。

这些**configurators**可以与`iris＃Addr，iris #Server，iris＃TLS，iris #AutheTLS，iris＃Listener` **final arguments（“hostConfigs”）**并排工作。

请注意，这些应用程序的主机**configurators**将与此应用程序可能创建的其他主机共享（使用`app.NewHost`），这意味着`app.NewHost`将在每次调用时执行这些**configurators**同样。

这些**configurators**应该在`app.Run`或`host.Serve / Listen`函数之前注册。

**[func (*Application) Logger](https://github.com/kataras/iris/blob/master/iris.go#L282)**

```go
func (app *Application) Logger() *golog.Logger
``` 

**Logger**返回在**app**中使用的**golog**记录器`instance(pointer)`。

可用级别： -**disable**- **fatal** - **error**- **warn** - **info** - **debug** 用法：`app.Logger()`。`SetLevel（“error”）`默认为**info**级别。

调用者可以使用应用程序的记录器，它是同一个`golog.Default` **logger**，也可以打印自定义日志。 用法：`app.Logger().Error/Errorf("...")`，`app.Logger().Warn/Warnf("...")` ， `app.Logger().Info/Infof("...")` ，`app.Logger().Debug/Debugf("...")`

设置一个或多个输出：`app.Logger()`。`SetOutput(io.Writer ...)`添加一个或多个输出：`app.Logger()`。`AddOutput(io.Writer ...)`

添加自定义级别需要导入`github.com/kataras/golog`包：

```go
//首先，我们创建一个golog.Level级别
//以便在日志功能中使用。
var SuccessLevel golog.Level = 6
//注册我们的级别，只有三个字段。
golog.Levels[SuccessLevel] = &golog.LevelMetadata{
	Name:    "success",
	RawText: "[SUCC]",
	// ColorfulText (Green Color[SUCC])
	ColorfulText: "\x1b[32m[SUCC]\x1b[0m",
}
``` 

示例: `app.Logger().SetLevel("success")` ，`app.Logger().Logf(SuccessLevel, "a custom leveled log message")`

**[func (*Application) NewHost](https://github.com/kataras/iris/blob/master/iris.go#L505)**

```go
func (app *Application) NewHost(srv *http.Server) *host.Supervisor
``` 

**NewHost**接受一个标准的`*http.Server`对象，完成该**srv**的必要缺失部分并返回一个新的，可立即使用 `host(supervisor)`。

**[func (*Application) RegisterView](https://github.com/kataras/iris/blob/master/iris.go#L310)**

```go
func (app *Application) RegisterView(viewEngine view.Engine)
``` 

**RegisterView**应该用于注册映射到根目录和模板文件扩展名的视图引擎。

**[func (*Application) Run](https://github.com/kataras/iris/blob/master/iris.go#L803)**

```go
func (app *Application) Run(serve Runner, withOrWithout ...Configurator) error
``` 

运行构建框架并启动所需的**Runner**，无论是否进行配置编辑。

每个**Application**实例只应调用一次**Run**，它会像**http.Server**一样阻塞。

如果需要在同一个**iris**实例上运行多个服务器，则创建一个新主机并通过`go NewHost(*http.Server).Serve/ListenAndServe`等手动运行它...或使用已创建的主机：`h： = NewHost(* http.Server)`运行`Raw(h.ListenAndServe)`，`WithCharset("UTF-8")`，`WithRemoteAddrHeader(“CF-Connecting-IP”)`

应用程序可以在以下运行程序的帮助下与任何类型的服务器或 **iris's host**：`Listener`，`Server`，`Addr`，`TLS`，`AutoTLS`和`Raw`。

**[func (*Application) SPA](https://github.com/kataras/iris/blob/master/iris.go#L474)**

```go
func (app *Application) SPA(assetHandler context.Handler) *router.SPABuilder
``` 

**SPA**接受**assetHandler**，它可以是`app.StaticHandler`或`app.StaticEmbeddedHandler`的结果。 当你想从**/index.html**导航到/自动它是一个辅助函数，它只是在**assetHandler**调用之前根据`IndexNames`和`AssetValidators`进行一些检查。

参阅：[: https://github.com/kataras/iris/tree/master/_examples/file-server/single-page-application](: https://github.com/kataras/iris/tree/master/_examples/file-server/single-page-application)

**[func (*Application) Shutdown](https://github.com/kataras/iris/blob/master/iris.go#L574)**

```go
func (app *Application) Shutdown(ctx stdContext.Context) error
``` 

**Shutdown**正常终止所有应用程序的服务器主机。 第一次失败时返回错误，否则返回nil。

**[func (*Application) SubdomainRedirect](https://github.com/kataras/iris/blob/master/iris.go#L219)**

```go
func (app *Application) SubdomainRedirect(from, to router.Party) router.Party
``` 

**SubdomainRedirect**注册一个路由器包装器，它在路由器尝试执行路由处理程序之前，尽可能快地将（子域）域（**StatusMovedPermanently**）重定向到另一个子域或根域。

它接收两个参数，它们是**from**和**to/target**位置，**from**也可以是一个通配符子域（`app.WildcardSubdomain()`）**to**由于显而易见的原因不允许成为通配符， **from** 能当**to**不是根域时，是`root(app)`，反之亦然。

用法： `www：= app.Subdomain(“www”)`< - 与`app.Party(“www.”)`相同`app.SubdomainRedirect(app，www)`这将重定向所有 **http://mydomain.com/%anypath％** 到 **http(s)://www.mydomain.com/%anypath%**。

一个或多个子域重定向可用于同一个应用实例。

如果您需要有关此实现的更多信息，那么您必须通过`core/router＃NewSubdomainRedirectWrapper`函数进行导航。

示例： [ https://github.com/kataras/iris/tree/master/_examples/subdomains/redirect]( https://github.com/kataras/iris/tree/master/_examples/subdomains/redirect)

**[func (*Application) View](https://github.com/kataras/iris/blob/master/iris.go#L323)**

```go
func (app *Application) View(writer io.Writer, filename string, layout string, bindingData interface{}) error
``` 

**View**执行并将模板文件的结果写入。

* 第一个参数是编写解析模板的编写器。 
* 第二个参数是相对于模板目录，模板文件名，包括扩展名。 
* 第三个参数是布局，可以是空字符串。 **Forth**参数是模板的可绑定数据，可以是**nil**。

使用`context.View`将模板呈现给客户端。 失败时返回错误，否则返回**nil**。

**[func (*Application) WWW](https://github.com/kataras/iris/blob/master/iris.go#L195)**

```go
func (app *Application) WWW() router.Party
``` 

**WWW**创建并返回`www.`子域。 与`app.Subdomain("www")`或`app.Party("www.")`的区别在于`app.WWW()`方法将路由器包装为所有**http(s)://mydomain.com** 将重定向到**http(s)://mydomain.com**。 其他子域名可以使用app注册：`sub：= app.Subdomain("mysubdomain")`，可以使用`www := app.WWW()`注册子子域名;www.Subdomain("wwwchildSubdomain")。

**[type Configuration](https://github.com/kataras/iris/blob/master/configuration.go#L362)**

```go
type Configuration struct {    
    //IgnoreServerErrors将导致忽略匹配的“错误”
    //来自主应用程序的`Run`函数。
    //这是一个字符串片段，而不是一段错误
    //用户可以使用yaml或toml配置文件注册这些错误
    //像其他配置字段一样。
    //更多参见`WithoutServerError（...）`函数。
    //
    // 示例: https://github.com/kataras/iris/tree/master/_examples/http-listening/listen-addr/omit-server-errors
    //
    //默认为空切片。
    // IgnoreServerErrors [] string`json：“ignoreServerErrors，omitempty”yaml：“IgnoreServerErrors”toml：“IgnoreServerErrors”`    
    // DisableStartupLog如果设置为true则关闭服务器启动时的写标题。
    //默认为false。
    DisableStartupLog bool `json:"disableStartupLog,omitempty" yaml:"DisableStartupLog" toml:"DisableStartupLog"`
    // DisableInterruptHandler如果设置为true，则会禁用自动正常服务器关闭
    //当按下control / cmd + C时
    //如果您计划通过自定义host.Task自行处理此问题，请将此项设置为true。
    //
    //默认为false。
    DisableInterruptHandler bool `json:"disableInterruptHandler,omitempty" yaml:"DisableInterruptHandler" toml:"DisableInterruptHandler"`

    // DisableVersionChecker如果为true，则不会通知进程任何可用的更新。
    // 
    // 默认为false。
    DisableVersionChecker bool `json:"disableVersionChecker,omitempty" yaml:"DisableVersionChecker" toml:"DisableVersionChecker"`
    
    // DisablePathCorrection更正并将请求的路径重定向到已注册的路径
    //例如，如果请求/home/ 但找不到此Route的处理程序，
    //然后路由器检查/home处理程序是否存在，如果是，
    //（permant）将客户端重定向到正确的路径/home
    //
    //默认为false。
    DisablePathCorrection bool `json:"disablePathCorrection,omitempty" yaml:"DisablePathCorrection" toml:"DisablePathCorrection"`

    //当EnablePathEscape为true时，它会转义路径，即命名参数（如果是 any 方法）。
    //如果你想要这样的东西https://github.com/kataras/iris/issues/135工作，改为fals
    //什么时候需要禁用（false）它：
    //接受斜杠'/'的参数
    //请求：http://localhost：8080/details/Project％2FDelta
    // ctx.Param（“project”）返回原始命名参数：Project％2FDelta
    //您可以使用net / url手动转义它：
    // projectName，_：= url.QueryUnescape（c.Param（“project”）。
    //
    //默认为false。
    EnablePathEscape bool `json:"enablePathEscape,omitempty" yaml:"EnablePathEscape" toml:"EnablePathEscape"`
    
    //当此字段为true时，启用EnableOptimization
    //然后应用程序尝试在可能的情况下优化以获得最佳性能。
    //
    //默认为false。
    EnableOptimizations bool `json:"enableOptimizations,omitempty" yaml:"EnableOptimizations" toml:"EnableOptimizations"`
    // FireMethodNotAllowed如果它是真正的路由器检查StatusMethodNotAllowed（405）和
    //触发405错误而不是404错误
    // 默认为false。
    FireMethodNotAllowed bool `json:"fireMethodNotAllowed,omitempty" yaml:"FireMethodNotAllowed" toml:"FireMethodNotAllowed"`

    // DisableBodyConsumptionOnUnmarshal管理上下文的正文阅读器/绑定器的阅读行为。
    //如果设置为true那么它
    //通过`context.UnmarshalBody / ReadJSON / ReadXML`禁用身体消耗。
    //
    //默认情况下，io.ReadAll`用于从`context.Request.Body中读取正文，它是一个`io.ReadCloser`，
    //如果此字段设置为true，则将创建一个新缓冲区以从请求正文中读取。
    //正文之前不会更改正文和现有数据
    // context.UnmarshalBody / ReadJSON / ReadXML将不被使用。
    DisableBodyConsumptionOnUnmarshal bool `json:"disableBodyConsumptionOnUnmarshal,omitempty" yaml:"DisableBodyConsumptionOnUnmarshal" toml:"DisableBodyConsumptionOnUnmarshal"`
    
    // DisableAutoFireStatusCode如果为true则关闭http错误状态码处理程序自动执行
    // from（`context.StatusCodeNotSuccessful`，默认为<200 ||> = 400）。
    //如果这是false，那么对于直接错误触发，然后手动调用“context＃FireStatusCode（statusCode）”。
    //
    //默认情况下，当调用“context.StatusCode（code）”时，将触发自定义http错误处理程序，
    //代码应该与`context.StatusCodeNotSuccessful`的结果相等，以便作为“http错误处理程序”接收。
    //
    //开发人员可能希望将此选项设置为true以便手动调用
    //需要时通过“context＃FireStatusCode（<200 ||> = 400）”处理错误处理程序。
    //正在通过app.OnErrorCode（code，handler）注册HTTP自定义错误处理程序。
    //
    //默认为false。
    DisableAutoFireStatusCode bool `json:"disableAutoFireStatusCode,omitempty" yaml:"DisableAutoFireStatusCode" toml:"DisableAutoFireStatusCode"`

    //任何类型的日期时间解析的TimeFormat时间格式
    //默认为“Mon，02 Jan 2006 15:04:05 GMT”。
    TimeFormat string `json:"timeFormat,omitempty" yaml:"TimeFormat" toml:"TimeFormat"`

    //用于各种渲染的字符集字符编码
    //用于模板和其他响应
    //默认为“UTF-8”。
    Charset string `json:"charset,omitempty" yaml:"Charset" toml:"Charset"`
    //PostMaxMemory设置最大发布数据大小
    //客户端可以发送到服务器，这是不同的
    //来自可以修改的overral请求正文大小
    //由`context＃SetMaxRequestBodySize`或`iris＃LimitRequestBodySize`。
    //
    //如果您愿意，可以默认为32MB或32 << 20。
    PostMaxMemory int64 `json:"postMaxMemory" yaml:"PostMaxMemory" toml:"PostMaxMemory"`
    //上下文值各种功能的键。
    //
    //TranslateLanguageContextKey和TranslateFunctionContextKey由i18n处理程序/中间件使用
    //目前我们只有一个：https：//github.com/kataras/iris/tree/master/middleware/i18n。
    //
    //默认为“iris.translate”和“iris.language”
    TranslateFunctionContextKey string `json:"translateFunctionContextKey,omitempty" yaml:"TranslateFunctionContextKey" toml:"TranslateFunctionContextKey"`
    //用于i18n的TranslateLanguageContextKey。
    //
    //默认为“iris.language”
    TranslateLanguageContextKey string `json:"translateLanguageContextKey,omitempty" yaml:"TranslateLanguageContextKey" toml:"TranslateLanguageContextKey"`

    //GetViewLayoutContextKey是上下文的用户值'键的关键
    //用于设置模板
    //来自中间件或主处理程序的布局。
    //覆盖父级或配置。
    //
    //默认为“iris.ViewLayout”
    ViewLayoutContextKey string `json:"viewLayoutContextKey,omitempty" yaml:"ViewLayoutContextKey" toml:"ViewLayoutContextKey"`
    // GetViewDataContextKey是上下文的用户值'键的关键
    //用于设置模板
    //绑定来自中间件或主处理程序的数据。
    //
    //默认为“iris.viewData”
    ViewDataContextKey string `json:"viewDataContextKey,omitempty" yaml:"ViewDataContextKey" toml:"ViewDataContextKey"`
    // RemoteAddrHeaders是允许的请求标头名称
    //可以有效解析客户端的IP。
    //默认情况下，没有“X-”标题可以安全地用于检索
    //客户端的IP地址，因为这些标头可以手动更改
    // 客户端。 但有时候例如在代理后面是有用的
    //你想启用“X-Forwarded-For”或cloudflare
    //你想要启用“CF-Connecting-IP”
    //可以允许`ctx.RemoteAddr（）`使用任何头
    //客户端可能发送的。
    
    //默认为空的map，但示例用法为：
    // RemoteAddrHeaders {
    //	"X-Real-Ip":             true,
    //  "X-Forwarded-For":       true,
    // 	"CF-Connecting-IP": 	 true,
    //	}
    //
   //查看`context.RemoteAddr（）`了解更多信息。
    RemoteAddrHeaders map[string]bool `json:"remoteAddrHeaders,omitempty" yaml:"RemoteAddrHeaders" toml:"RemoteAddrHeaders"`

    //其他是自定义的动态选项，可以为空。
    //此字段仅供您用于设置所需的任何应用程序选项。
    //
    //默认为非零空的map。
    Other map[string]interface{} `json:"other,omitempty" yaml:"Other" toml:"Other"`
    //包含已过滤或未导出的字段
}
``` 

配置**iris**实例的整个配置，这些也可以通过选项传递，查看此文件的顶部（**configuration.go**）。 配置是有效的**OptionSetter**。

**[func DefaultConfiguration](https://github.com/kataras/iris/blob/master/configuration.go#L757)**

```go
func DefaultConfiguration() Configuration
``` 

**DefaultConfiguration**返回**iris**的默认配置，填充主配置

**[func TOML](https://github.com/kataras/iris/blob/master/configuration.go#L123)**

```go
func TOML(filename string) Configuration
``` 

TOML从兼容toml的文档文件中读取配置。 有关toml实现的更多信息，请访问：[https：//github.com/toml-lang/toml](https：//github.com/toml-lang/toml)

接受配置文件的绝对路径。 将通过出现错误消息的恐慌向用户显示错误。 文件不存在或格式不正确时可能会出错。

注意：如果**char'〜'**作为**filename**传递，那么它会尝试从 `$home_directory + iris.tml`加载并返回配置，有关详细信息，请参阅`WithGlobalConfiguration`。

用法：`app.Configure(iris.WithConfiguration(iris.YAML(“myconfig.tml”)))`或`app.Run([iris.Runner]，iris.WithConfiguration(iris.YAML(“myconfig.tml”)))`。

**[func YAML](https://github.com/kataras/iris/blob/master/configuration.go#L89)**

```go
func YAML(filename string) Configuration
``` 

**YAML**从`configuration.yml`文件中读取配置。

接受**cfg.yml**的绝对路径。 将通过出现错误消息的恐慌向用户显示错误。 **cfg.yml**不存在或格式不正确时可能会出错。

注意：如果**char'〜'**作为**filename**传递，那么它会尝试从`$home_directory + iris.yml`加载并返回配置，有关详细信息，请参阅`WithGlobalConfiguration`。

用法: `app.Configure(iris.WithConfiguration(iris.YAML("myconfig.yml")))` 或 `app.Run([iris.Runner], iris.WithConfiguration(iris.YAML("myconfig.yml")))`.

**[func (Configuration) GetCharset](https://github.com/kataras/iris/blob/master/configuration.go#L588)**

```go
func (c Configuration) GetCharset() string
``` 

**GetCharset**返回`Configuration＃Charset`，用于模板和其余响应的各种呈现的字符编码。

**[func (Configuration) GetDisableAutoFireStatusCode](https://github.com/kataras/iris/blob/master/configuration.go#L575)**

```go
func (c Configuration) GetDisableAutoFireStatusCode() bool
``` 

**GetDisableAutoFireStatusCode**返回`Configuration＃DisableAutoFireStatusCode`。 当http错误状态代码处理程序自动执行关闭时，返回true。

**[func (Configuration) GetDisableBodyConsumptionOnUnmarshal](https://github.com/kataras/iris/blob/master/configuration.go#L569)**

```go
func (c Configuration) GetDisableBodyConsumptionOnUnmarshal() bool
``` 

**GetDisableBodyConsumptionOnUnmarshal**返回`Configuration＃GetDisableBodyConsumptionOnUnmarshal`，管理上下文的正文阅读器/绑定器的阅读行为。 如果返回true，则禁用`context.UnmarshalBody / ReadJSON / ReadXML`的正文消耗。

GetDisableBodyConsumptionOnUnmarshal returns the Configuration#GetDisableBodyConsumptionOnUnmarshal, manages the reading behavior of the context's body readers/binders. If returns true then the body consumption by the `context.UnmarshalBody/ReadJSON/ReadXML` is disabled.

默认情况下，**io.ReadAll**用于从`context.Request.Body`读取请求体，它是一个`io.ReadCloser`，如果此字段设置为 **true**，则将创建一个新缓冲区以从请求体中读取。 在上下文之前不会更改正文和现有数据.`UnmarshalBody / ReadJSON / ReadXML`将不会被使用。

**[func (Configuration) GetDisablePathCorrection](https://github.com/kataras/iris/blob/master/configuration.go#L539)**

```go
func (c Configuration) GetDisablePathCorrection() bool
``` 

**GetDisablePathCorrection**返回`Configuration＃DisablePathCorrection`，`DisablePathCorrection`更正并将请求的路径重定向到已注册的路径，例如，如果请求 **/home/path**但没有找到此Route的处理程序，则路由器检查 **/home**处理程序是否存在，如果是， **(permant)** 将客户端重定向到正确的路径 **/home**。

**[func (Configuration) GetEnableOptimizations](https://github.com/kataras/iris/blob/master/configuration.go#L551)**

```go
func (c Configuration) GetEnableOptimizations() bool
``` 

**GetEnableOptimizations**返回应用程序是否启用了性能优化。

**[func (Configuration) GetEnablePathEscape](https://github.com/kataras/iris/blob/master/configuration.go#L545)**

```go
func (c Configuration) GetEnablePathEscape() bool
``` 

**GetEnablePathEscape**是`Configuration＃EnablePathEscape`，当它转义路径时返回true，即命名参数（如果有的话）。

**[func (Configuration) GetFireMethodNotAllowed](https://github.com/kataras/iris/blob/master/configuration.go#L556)**

```go
func (c Configuration) GetFireMethodNotAllowed() bool
``` 

**GetFireMethodNotAllowed**返回`Configuration＃FireMethodNotAllowed`。

**[func (Configuration) GetOther](https://github.com/kataras/iris/blob/master/configuration.go#L652)**

```go
func (c Configuration) GetOther() map[string]interface{}
``` 

**GetOther**返回`Configuration＃Other`映射。

**[func (Configuration) GetPostMaxMemory](https://github.com/kataras/iris/blob/master/configuration.go#L598)**

```go
func (c Configuration) GetPostMaxMemory() int64
``` 

**GetPostMaxMemory**返回客户端可以发送到服务器的最大配置的发布数据大小，这与可以通过`context＃SetMaxRequestBodySize`或`iris＃LimitRequestBodySize`修改的**overral**请求主体大小不同。

如果您愿意，默认为32MB或32 << 20。

**[func (Configuration) GetRemoteAddrHeaders](https://github.com/kataras/iris/blob/master/configuration.go#L647)**

```go
func (c Configuration) GetRemoteAddrHeaders() map[string]bool
``` 

**GetRemoteAddrHeaders**返回允许的请求标头名称，这些名称可以有效解析客户端的**IP**。 默认情况下，没有“X-”标头可以安全地用于检索客户端的**IP**地址，因为这些标头可以由客户端手动更改。
但有时候很有用，例如，在代理服务器后面你想启用**X-Forwarded-For**或者当你想要启用**CF-Connecting-IP**的**cloudflare**时，你可以允许`ctx.RemoteAddr()` 使用客户端可能发送的任何标头。

默认为空map，但示例用法为：**RemoteAddrHeaders**

```go
    {
    "X-Real-Ip":             true,
    "X-Forwarded-For":       true,
    "CF-Connecting-IP": 	 true,
    }
``` 

看看`context.RemoteAddr()`了解更多。

**[func (Configuration) GetTimeFormat](https://github.com/kataras/iris/blob/master/configuration.go#L581)**

```go
func (c Configuration) GetTimeFormat() string
``` 

**GetTimeFormat**为任何类型的日期时间解析返回`Configuration＃TimeFormat`格式。

**[func (Configuration) GetTranslateFunctionContextKey](https://github.com/kataras/iris/blob/master/configuration.go#L604)**

```go
func (c Configuration) GetTranslateFunctionContextKey() string
``` 

**GetTranslateFunctionContextKey**返回配置的`TranslateFunctionContextKey`值，用于**i18n**。

**[func (Configuration) GetTranslateLanguageContextKey](https://github.com/kataras/iris/blob/master/configuration.go#L610)**

```go
func (c Configuration) GetTranslateLanguageContextKey() string
``` 

**GetTranslateLanguageContextKey**返回配置的`TranslateLanguageContextKey`值，用于**i18n**。

**[func (Configuration) GetVHost](https://github.com/kataras/iris/blob/master/configuration.go#L530)**

```go
func (c Configuration) GetVHost() string
``` 

**GetVHost**返回未导出的**vhost**配置字段。

如果原始地址以：**443**或 **：80** 结束，它将返回没有端口的主机。 如果原始地址是：**https**或：**http**，它将返回**localhost**。 如果原始地址为**0.0.0.0**，则返回**localhost**。

**[func (Configuration) GetViewDataContextKey](https://github.com/kataras/iris/blob/master/configuration.go#L625)**

```go
func (c Configuration) GetViewDataContextKey() string
``` 

**GetViewDataContextKey**返回上下文的用户值'键 的键，用于从中间件或主处理程序设置模板绑定数据。

**[func (Configuration) GetViewLayoutContextKey](https://github.com/kataras/iris/blob/master/configuration.go#L618)**

```go
func (c Configuration) GetViewLayoutContextKey() string
``` 

**GetViewLayoutContextKey**返回上下文的用户值'键 的键，用于从中间件或主处理程序设置模板布局。 覆盖父级或配置。


#### [type Configurator](https://github.com/kataras/iris/blob/master/configuration.go#L166)

```go
type Configurator func(*Application)
``` 

**Configurator**只是一个接受框架实例的接口。

它可用于通过`Configure`注册自定义配置，以便修改框架实例。

目前，**Configurator**用于描述配置的字段值。

**[func WithCharset) GetViewLayoutContextKey](https://github.com/kataras/iris/blob/master/configuration.go#L287)**

```go
func WithCharset(charset string) Configurator
``` 

**WithCharset**设置**Charset**设置。

参见`Configuration`。

**[func WithConfiguration](https://github.com/kataras/iris/blob/master/configuration.go#L664)**

```go
func WithConfiguration(c Configuration) Configurator
``` 

**WithConfiguration** 将 **c** 值设置为框架的配置。

用法：`app.Run(iris.Addr(“：8080”)，iris.WithConfiguration(iris.Configuration {/ * fields here * /}))`或``或`iris.WithConfiguration(iris.TOML(”./ cfg / iris.tml“))`iris.WithConfiguration(iris.YAML(“./cfg/iris.yml"))

**[func WithOtherValue](https://github.com/kataras/iris/blob/master/configuration.go#L350)**

```go
func WithOtherValue(key string, val interface{}) Configurator
``` 

**WithOtherValue** 根据“其他”设置的键添加值。

参见`Configuration`。

**[func WithPostMaxMemory](https://github.com/kataras/iris/blob/master/configuration.go#L299)**

```go
func WithPostMaxMemory(limit int64) Configurator
``` 

**WithPostMaxMemory**设置客户端可以发送到服务器的最大发布数据大小，这不同于可以通过`context＃SetMaxRequestBodySize`或`iris＃LimitRequestBodySize`修改的**overral**请求主体大小。

如果您愿意，默认为32MB或32 << 20。

**[func WithRemoteAddrHeader](https://github.com/kataras/iris/blob/master/configuration.go#L320)**

```go
func WithRemoteAddrHeader(headerName string) Configurator
``` 

**WithRemoteAddrHeader**启用或添加可用于验证客户端真实**IP**的新或现有请求标头名称。

默认情况下，没有 **X-** 标头可以安全地用于检索客户端的 **IP** 地址，因为这些标头可以由客户端手动更改。 
但有时候很有用，例如，在代理服务器后面你想启用**X-Forwarded-For**或者当你想要启用**CF-Connecting-IP**的**cloudflare**时，你可以允许`ctx.RemoteAddr()` 使用客户端可能发送的任何标头。

默认为空地图，但示例用法为：`WithRemoteAddrHeader(“X-Forwarded-For”)`

看看`context.RemoteAddr()`了解更多。

**[func WithTimeFormat](https://github.com/kataras/iris/blob/master/configuration.go#L278)**

```go
func WithTimeFormat(timeformat string) Configurator
``` 

使用时间格式设置时间格式设置。

参见`Configuration`。

**[func WithoutRemoteAddrHeader](https://github.com/kataras/iris/blob/master/configuration.go#L338)**

```go
func WithoutRemoteAddrHeader(headerName string) Configurator
``` 

**WithoutRemoteAddrHeader**禁用可用于验证和解析客户端的真实**IP**的现有请求标头名称。

请注意，**RemoteAddrHeaders**已默认为空映射，因此如果您之前未通过配置或通过**WithRemoteAddrHeader**添加允许的标头，则无需调用此配置器。

看看`context.RemoteAddr()`了解更多。

**[func WithoutServerError](https://github.com/kataras/iris/blob/master/configuration.go#L194)**

```go
func WithoutServerError(errors ...error) Configurator
``` 

**WithoutServerError**将导致忽略主应用程序的`Run`函数中匹配的**errors**。

用法：`errors := app.Run(iris.Addr(“：8080”)，iris.WithoutServerError(iris.ErrServerClosed))`如果服务器的错误是`http/iris＃ErrServerClosed`，将返回`nil`。

请参阅`Configuration＃IgnoreServerErrors [] string`。

示例： [Example: https://github.com/kataras/iris/tree/master/_examples/http-listening/listen-addr/omit-server-errors](Example: https://github.com/kataras/iris/tree/master/_examples/http-listening/listen-addr/omit-server-errors)

#### [type Context](https://github.com/kataras/iris/blob/master/go19.go#L19)

```go
type Context = context.Context
``` 

**Context**是客户端的**midle-man**服务器的“对象”。

正在从每个连接上的`sync.Pool`获取新上下文。 **Context**是`iris http`流中最重要的东西。

开发人员通过**Context**发送对客户端请求的响应。 开发人员通过**Context**从客户端的请求中获取请求信息。

#### [type CookieOption](https://github.com/kataras/iris/blob/master/go19.go#L83)

```go
type CookieOption = context.CookieOption
``` 

**CookieOption**是上下文的方法接受的函数类型，如`SetCookieKV`，`RemoveCookie`和`SetCookie`作为修改结束cookie形式的（最后）可变参数输入参数。

任何自定义或构建的**CookieOption**都是有效的，请参阅`CookiePath`，`CookieCleanPath`，`CookieExpires`和`CookieHTTPOnly`了解更多信息。

`context/Context＃CookieOption`的别名。

#### [type ExecutionOptions](https://github.com/kataras/iris/blob/master/go19.go#L73)

```go
type ExecutionOptions = router.ExecutionOptions
``` 

**ExecutionOptions**是一组默认行为，可以更改这些行为，以便轻松自定义路径处理程序的执行流程。

有关更多信息，请参阅`ExecutionRules`和`core/router/Party＃SetExecutionRules`。

#### [type ExecutionRules](https://github.com/kataras/iris/blob/master/go19.go#L69)

```go
type ExecutionRules = router.ExecutionRules
``` 

**ExecutionRules**控制处理程序本身之外的路由处理程序的执行。 用法`：Party＃SetExecutionRules（ExecutionRules){Done: ExecutionOptions{Force: true},})`

有关更多信息，请参阅`core/router/Party＃SetExecutionRules`。 

例： [https://github.com/kataras/iris/tree/master/_examples/mvc/middleware/without-ctx-next](https://github.com/kataras/iris/tree/master/_examples/mvc/middleware/without-ctx-next)

#### [type Handler](https://github.com/kataras/iris/blob/master/go19.go#L44)

```go
type Handler = context.Handler
``` 

处理程序响应**HTTP**请求。它将回复标头和数据写入`Context.ResponseWriter()`然后返回。 返回请求完成的信号; 在完成**Handler**调用之后或同时使用**Context**是无效的。

根据**HTTP**客户端软件，**HTTP**协议版本以及客户端和**iris**服务器之间的任何中介，可能无法在写入`context.ResponseWriter()`之后从`Context.Request()`。**Body**中读取。 谨慎的处理程序应该首先阅读`Context.Request()`。**Body**，然后回复。

除了阅读正文外，处理程序不应修改提供的**Context**。

如果**Handler**发生混乱，服务器（**Handler**的调用者）会认为恐慌的影响与活动请求隔离。 它恢复了异常，将堆栈跟踪记录到服务器错误日志，并挂断连接。

#### [type Map](https://github.com/kataras/iris/blob/master/go19.go#L46)

```go
type Map = context.Map
``` 

`Map`是`map [string] interface {}`的快捷方式。

#### [type Party](https://github.com/kataras/iris/blob/master/go19.go#L59)

```go
type Party = router.Party
``` 

**Party**只是具有相同前缀并共享相同中间件的路由的组合连接器。 派对也可以被命名为“加入”或“节点”或“群组”，派对因为它很有趣而被选中。

看看`core/router＃APIBuilder`的实现。

`core/router＃Party`的快捷方式，在使用`PartyFunc`时很有用。

#### [type Runner](https://github.com/kataras/iris/blob/master/iris.go#L594)

```go
type Runner func(*Application) error
``` 

**Runner**只是一个接受框架实例并返回错误的接口。

它可用于向`Run`注册自定义运行器，以便设置框架的服务器侦听操作。

当前**Runner**用于声明内置的服务器侦听器。

有关更多信息，请参阅`Run`。

**[func Addr](https://github.com/kataras/iris/blob/master/iris.go#L655)**

```go
func Addr(addr string, hostConfigs ...host.Configurator) Runner
``` 

**Addr**可以用作`Run`方法的参数。 它接受用于构建服务器的主机地址和侦听该主机和端口的侦听器。

**Addr**的形式应为`[host]：port`，即`localhost：8080`或`：8080`。

第二个参数是可选的，它接受一个或多个`func(* host.Configurator)`，它们正在该特定主机上执行，该函数将创建该主机以启动服务器。 
通过主机配置程序，您可以配置后端主机管理程序，即添加关闭，服务或错误的事件。 可以在以下网址找到此用例的示例：[https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go](https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go)也可以查看`ConfigureHost`。

有关更多信息，请参阅`Run`。

**[func AutoTLS](https://github.com/kataras/iris/blob/master/iris.go#L723)**

```go
func AutoTLS(
    addr string,
    domain string, email string,
    hostConfigs ...host.Configurator) Runner
    }
``` 

**AutoTLS**可以用作`Run`方法的参数。 它将使用**autocert** `golang/x`包即时创建的认证启动应用程序的安全服务器，因此**localhost**可能无法正常工作，在**production**机器上使用它。

**Addr**的形式应为`[host]：port`，即**mydomain.com:443**。

白名单域由**domain**参数中的空格分隔，即**iris-go.com**，可以与**addr**不同。 如果为空，则当前允许所有主机。
建议不要这样做，因为它会打开潜在的攻击，客户端通过**IP**地址连接到服务器并假装要求输入错误的主机名。
**Manager**将尝试错误地获取该主机的证书，最终达到**CA**的证书请求速率限制，并且无法获得实际证书。

对于使用非公开的**e-mail**，**letsencrypt**需要为您自己的安全。

注意：`AutoTLS`将为您启动一个新服务器，它将所有**http**版本重定向到他们的**https**，包括子域。

最后一个参数是可选的，它接受一个或多个`func(* host.Configurator)`，它们正在该特定主机上执行，该函数将创建该主机以启动服务器。
通过主机配置程序，您可以配置后端主机管理程序，即添加关闭，服务或错误的事件。
可以在以下网址找到此用例的示例：[https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go](https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go)也可以查看`ConfigureHost`。

用法：app.Run(iris.AutoTLS(“iris-go.com:443”，“iris-go.com www.iris-go.com”，“mail@example.com”))

有关更多信息，请参阅`Run`和`core/host/Supervisor＃ListenAndServeAutoTLS`。

**[func Listener](https://github.com/kataras/iris/blob/master/iris.go#L609)**

```go
func Listener(l net.Listener, hostConfigs ...host.Configurator) Runner
``` 

**Listener**可以用作`Run`方法的参数。 它可以通过服务器的**Serve**启动带有自定义`net.Listener`的服务器。

第二个参数是可选的，它接受一个或多个`func(* host.Configurator)`，它们正在该特定主机上执行，该函数将创建该主机以启动服务器。
通过主机配置程序，您可以配置后端主机管理程序，即添加关闭，服务或错误的事件。
可以在以下网址找到此用例的示例：[https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go](https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go)也可以查看`ConfigureHost`。

有关更多信息，请参阅`Run`。

**[func Raw](https://github.com/kataras/iris/blob/master/iris.go#L743)**

```go
func Raw(f func() error) Runner
``` 

Raw可以用作`Run`方法的参数。 它接受任何 **(listen)** 函数返回错误，此函数应该是阻塞并仅在服务器退出或导致致命错误时返回错误。

使用此选项，您不仅限于 **iris** 可以默认运行的服务器。

有关更多信息，请参阅`Run`。

**[func Server](https://github.com/kataras/iris/blob/master/iris.go#L631)**

```go
func Server(srv *http.Server, hostConfigs ...host.Configurator) Runner
``` 

服务器可以用作`Run`方法的参数。 它可以使用`*http.Server`启动服务器。

第二个参数是可选的，它接受一个或多个`func(* host.Configurator)`，它们正在该特定主机上执行，该函数将创建该主机以启动服务器。
通过主机配置程序，您可以配置后端主机管理程序，即添加关闭，服务或错误的事件。 
可以在以下网址找到此用例的示例：[https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go](https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go)也可以查看`ConfigureHost`。

有关更多信息，请参阅`Run`。

**[func TLS](https://github.com/kataras/iris/blob/master/iris.go#L681)**

```go
func TLS(addr string, certFile, keyFile string, hostConfigs ...host.Configurator) Runner
``` 

**TLS**可以用作`Run`方法的参数。 它将启动**Application**的安全服务器。

像以前一样使用它来使用`http.ListenAndServeTLS`函数。

**Addr**的形式应为`[host]：port`，即**localhost：443**或：**443**。 **CertFile**和**KeyFile**应该是带有扩展名的文件名。

第二个参数是可选的，它接受一个或多个`func(* host.Configurator)`，它们正在该特定主机上执行，该函数将创建该主机以启动服务器。
通过主机配置程序，您可以配置后端主机管理程序，即添加关闭，服务或错误的事件。
可以在以下网址找到此用例的示例：[https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go](https：//github.com/kataras/iris/blob/master/_examples/http-listening/notify-on-shutdown/main.go)也可以查看`ConfigureHost`。

有关更多信息，请参阅`Run`。

#### [type Supervisor](https://github.com/kataras/iris/blob/master/go19.go#L51)

```go
type Supervisor = host.Supervisor
``` 

**Supervisor** 快捷方式 是`host＃Supervisor`类型。 用于在常用**Runners**上添加**supervisor**配置程序，而无需导入`core/host`包。

#### [type UnmarshalerFunc](https://github.com/kataras/iris/blob/master/go19.go#L29)

```go
type UnmarshalerFunc = context.UnmarshalerFunc
``` 

**UnmarshalerFunc**是一个快捷方式，是`context＃UnmarshalerFunc`类型的别名，它实现了`context＃Unmarshaler`接口，用于通过自定义解码器读取请求的主体，其中大多数已经实现了像`json.Unmarshal，xml`，`context＃UnmarshalerFunc`。
**Unmarshal**，`yaml.Unmarshal`以及遵循最佳实践并符合**Go**标准库。

有关更多信息，请参阅`context＃UnmarshalBody`。

示例：[https：//github.com/kataras/iris/blob/master/_examples/http_request/read-custom-via-unmarshaler/main.go](https：//github.com/kataras/iris/blob/master/_examples/http_request/read-custom-via-unmarshaler/main.go)

#### 目录

| 路径 | 概要 |
| -----------|-------------|
| [cache](https://godoc.org/github.com/kataras/iris/cache) |       |
| [cache/cfg](https://godoc.org/github.com/kataras/iris/cache/cfg) |       |
| [cache/client](https://godoc.org/github.com/kataras/iris/cache/client) |       |
| [cache/client/rule](https://godoc.org/github.com/kataras/iris/cache/client/rule) |       |
| [cache/entry](https://godoc.org/github.com/kataras/iris/cache/entry) |       |
| [cache/ruleset](https://godoc.org/github.com/kataras/iris/cache/ruleset) | 包规则集提供了规则扩展的基本规则。      |
| [cache/uri](https://godoc.org/github.com/kataras/iris/cache/uri) |       |
| [context](https://godoc.org/github.com/kataras/iris/context) |       |
| [core/errors](https://godoc.org/github.com/kataras/iris/core/errors) |       |
| [core/handlerconv](https://godoc.org/github.com/kataras/iris/core/handlerconv) |       |
| [core/host](https://godoc.org/github.com/kataras/iris/core/host) |       |
| [core/maintenance](https://godoc.org/github.com/kataras/iris/core/maintenance) |       |
| [core/maintenance/version](https://godoc.org/github.com/kataras/iris/core/maintenance/version) |       |
| [core/memstore](https://godoc.org/github.com/kataras/iris/core/memstore) | 包memstore包含一个存储器，它只是具有不变能力的键值条目的集合。      |
| [core/netutil](https://godoc.org/github.com/kataras/iris/core/netutil) |       |
| [core/router](https://godoc.org/github.com/kataras/iris/core/router) |       |
| [core/router/macro](https://godoc.org/github.com/kataras/iris/core/router/macro) |       |
| [core/router/macro/interpreter/ast](https://godoc.org/github.com/kataras/iris/core/router/macro/interpreter/ast) |       |
| [core/router/macro/interpreter/lexer](https://godoc.org/github.com/kataras/iris/core/router/macro/interpreter/lexer) |       |
| [core/router/macro/interpreter/parser](https://godoc.org/github.com/kataras/iris/core/router/macro/interpreter/parser) |       |
| [core/router/macro/interpreter/token](https://godoc.org/github.com/kataras/iris/core/router/macro/interpreter/token) |       |
| [core/router/node](https://godoc.org/github.com/kataras/iris/core/router/node) |       |
| [hero](https://godoc.org/github.com/kataras/iris/hero) |       |
| [hero/di](https://godoc.org/github.com/kataras/iris/hero/di) | Package di为Iris Hero和Iris MVC新功能提供依赖注入。 |
| [httptest](https://godoc.org/github.com/kataras/iris/httptest) |       |
| [middleware/basicauth](https://godoc.org/github.com/kataras/iris/middleware/basicauth) | Package basicauth通过中间件提供http基本身份验证。|
| [middleware/i18n](https://godoc.org/github.com/kataras/iris/middleware/i18n) |包i18n通过中间件提供内部化和本地化。 |
| [middleware/logger](https://godoc.org/github.com/kataras/iris/middleware/logger) | 包记录器通过中间件提供请求记录。|
| [middleware/pprof](https://godoc.org/github.com/kataras/iris/middleware/pprof) |包pprof通过中间件提供本机pprof支持。       |
| [middleware/recaptcha](https://godoc.org/github.com/kataras/iris/middleware/recaptcha) |       |
| [middleware/recover](https://godoc.org/github.com/kataras/iris/middleware/recover) | 包恢复通过中间件为特定路由或整个应用程序提供恢复。      |
| [mvc](https://godoc.org/github.com/kataras/iris/mvc) |       |
| [sessions](https://godoc.org/github.com/kataras/iris/sessions) |       |
| [sessions/sessiondb/badger](https://godoc.org/github.com/kataras/iris/sessions/sessiondb/badger) |       |
| [sessions/sessiondb/boltdb](https://godoc.org/github.com/kataras/iris/sessions/sessiondb/boltdb) |       |
| [sessions/sessiondb/redis](https://godoc.org/github.com/kataras/iris/sessions/sessiondb/redis) |       |
| [sessions/sessiondb/redis/service](https://godoc.org/github.com/kataras/iris/sessions/sessiondb/redis/service) |       |
| [typescript](https://godoc.org/github.com/kataras/iris/typescript) |Package typescript为typescript编译器提供了热重写器和可选的基于云的编辑器，称为“alm-tools”。       |
| [typescript/editor](https://godoc.org/github.com/kataras/iris/typescript/editor) |       |
| [typescript/npm](https://godoc.org/github.com/kataras/iris/typescript/npm) |       |
| [view](https://godoc.org/github.com/kataras/iris/view) |       |
| [websocket](https://godoc.org/github.com/kataras/iris/websocket) | 包websocket为iris框架提供了丰富的websocket支持。      |

**iris**导入26个包（图）并由561个包导入。

#### 翻译有错误或者不准确请多多包含，也请不吝指教，指出不足。是我们更好的完善，为社区做贡献。 `ZRothschild`