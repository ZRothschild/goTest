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
