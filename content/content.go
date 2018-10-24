package main

import (
	"time"
	"context"
	"fmt"
)

var key string="name"
var keys string="names"

func main() {
	ctx, cancel := context.WithDeadline(context.Background(),time.Now().Add(1*time.Second))
	//附加值
	valueCtx:=context.WithValue(ctx,key,"【监控1】")
	go watch(valueCtx)
	aa := context.WithValue(valueCtx,keys,"【监控3】")
	go watch(aa)
	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	fmt.Println(ctx.Err())
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}

func watch(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			//取出值
			fmt.Println(ctx.Value(key),"监控退出，停止了...")
			return
		default:
			//取出值
			fmt.Println(ctx.Value(key),"goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}