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