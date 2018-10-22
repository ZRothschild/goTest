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