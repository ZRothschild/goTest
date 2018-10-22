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