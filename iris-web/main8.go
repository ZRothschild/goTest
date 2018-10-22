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