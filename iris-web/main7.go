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
}
