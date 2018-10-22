package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()

	// - 标准 html      | iris.HTML(...)
	// - django         | iris.Django(...)
	// - pug(jade)      | iris.Pug(...)
	// - handlebars     | iris.Handlebars(...)
	// - amber          | iris.Amber(...)

	tmpl := iris.HTML("./web/views/main10", ".html")
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