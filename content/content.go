package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second*3)
	fmt.Println(ctx.Err())
	defer cancel() // 防止任务比超时时间短导致资源未释放
	// 启动协程
	go task(ctx)
	// 主协程需要等待，否则直接退出
	time.Sleep(time.Second * 20)
}

func task(ctx context.Context) {
	ch := make(chan struct{}, 0)
	// 真正的任务协程
	go func() {
		// 模拟10秒耗时任务 这里修改成大于4秒
		time.Sleep(time.Second * 5)
		ch <- struct{}{}
	}()
	select {
	case <-ch:
		fmt.Println("done")
	case <-ctx.Done():
		fmt.Println("timeout")
	}
}

// package main
//
// import (
// 	"context"
// 	"fmt"
// 	"time"
// )
//
// var key string = "name"
// var keys string = "names"
//
// func main() {
// 	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(1*time.Second))
// 	//附加值
// 	valueCtx := context.WithValue(ctx, key, "【监控1】")
// 	go watch(valueCtx)
// 	aa := context.WithValue(valueCtx, keys, "【监控3】")
// 	go watch(aa)
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("可以了，通知监控停止")
// 	fmt.Println(ctx.Err())
// 	cancel()
// 	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
// 	time.Sleep(5 * time.Second)
// }
//
// func watch(ctx context.Context) {
// 	for {
// 		select {
// 		case <-ctx.Done():
// 			//取出值
// 			fmt.Println(ctx.Value(key), "监控退出，停止了...")
// 			return
// 		default:
// 			//取出值
// 			fmt.Println(ctx.Value(key), "goroutine监控中...")
// 			time.Sleep(2 * time.Second)
// 		}
// 	}
// }
