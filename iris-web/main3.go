package main

import (
	stdContext "context"
	"time"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()

	iris.RegisterOnInterrupt(func() {
		timeout := 20 * time.Second
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