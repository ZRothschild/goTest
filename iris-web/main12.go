package main

import (
	"time"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
	"github.com/kataras/iris/sessions/sessiondb/redis"
	"github.com/kataras/iris/sessions/sessiondb/redis/service"
)

//使用redis 3.0.503版测试。
//对于Windows，请参阅：https：//github.com/ServiceStack/redis-windows
func main() {
	//替换正在运行的redis'服务器设置：
	db := redis.New(service.Config{
		Network:     service.DefaultRedisNetwork,
		Addr:        service.DefaultRedisAddr,
		Password:    "",
		Database:    "",
		MaxIdle:     0,
		MaxActive:   0,
		IdleTimeout: service.DefaultRedisIdleTimeout,
		Prefix:      ""}) //可选择在redis服务器之间配置网桥

	//当control + C / cmd + C时关闭连接
	iris.RegisterOnInterrupt(func() {
		db.Close()
	})

	defer db.Close() //如果应用程序出错，请关闭数据库连接。

	sess := sessions.New(sessions.Config{
		Cookie:  "sessionscookieid",
		Expires: 45 * time.Minute}, // <= 0意味着不会过期。 默认为0。
	)

	// 重要：
	sess.UseDatabase(db)

	//其余代码保持不变
	app := iris.New()

	app.Get("/", func(ctx iris.Context) {
		ctx.Writef("You should navigate to the /set, /get, /delete, /clear,/destroy instead")
	})
	app.Get("/set", func(ctx iris.Context) {
		s := sess.Start(ctx)
		//设置会话值
		s.Set("name", "iris")

		//测试获取
		ctx.Writef("All ok session value of the 'name' is: %s", s.GetString("name"))
	})

	app.Get("/set/{key}/{value}", func(ctx iris.Context) {
		key, value := ctx.Params().Get("key"), ctx.Params().Get("value")
		s := sess.Start(ctx)
		//设置会话值
		s.Set(key, value)

		//测试获取
		ctx.Writef("All ok session value of the '%s' is: %s", key, s.GetString(key))
	})

	app.Get("/get", func(ctx iris.Context) {
		//获取一个特定的键，作为字符串，如果没有找到只返回一个空字符串
		name := sess.Start(ctx).GetString("name")

		ctx.Writef("The 'name' on the /set was: %s", name)
	})

	app.Get("/get/{key}", func(ctx iris.Context) {
		//获取一个特定的键，作为字符串，如果没有找到只返回一个空字符串
		name := sess.Start(ctx).GetString(ctx.Params().Get("key"))

		ctx.Writef("The name on the /set was: %s", name)
	})

	app.Get("/delete", func(ctx iris.Context) {
		//删除特定 key
		sess.Start(ctx).Delete("name")
	})

	app.Get("/clear", func(ctx iris.Context) {
		//删除所有条目
		sess.Start(ctx).Clear()
	})

	app.Get("/destroy", func(ctx iris.Context) {
		// 销毁，删除整个会话数据和cookie
		sess.Destroy(ctx)
	})

	app.Get("/update", func(ctx iris.Context) {
		//使用新日期更新过期日期
		sess.ShiftExpiration(ctx)
	})

	app.Run(iris.Addr(":8080"), iris.WithoutServerError(iris.ErrServerClosed))
}