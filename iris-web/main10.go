package main

import "github.com/kataras/iris"

func main() {
	app := iris.New()
	// $ go get -u github.com/shuLhan/go-bindata/...  切记安装时不要忘记三个点 ...
	// $ go-bindata ./web/views/main10/...
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