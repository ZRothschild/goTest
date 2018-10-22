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


	// More than one route can contain the same path with a different http mapped method.
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

	// GET: http://localhost:8080/admin/login
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