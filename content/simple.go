package main

import (
	"time"
	"context"
	"fmt"
)


/**
目前总结两种使用场景：
1. 主动停止goroutine
2. 传递数据

Context是什么？
Context直译为上下文，我们来看一下Context包中对它的说明，在go1.7之后已经添加到了标准库中，我们之间可以在/src/context中查看。

1. 包context定义了Context类型，并包含deadlines(结束时间),cancelation signals（取消信号），和其他的请求api范围的值。
2. 对服务器的请求应该创建一个Context，服务器发出的外向请求应该接受Context。
链式的函数调用之间必须传递Context，随意的更换它使用一个导出Context使用WithCancel、WithDeadline、WithTimeout、WithValue。
当一个Context取消所有从它导出的Contexts都会被取消。
3. WithCancel、WithDeadline、WithTimeout函数携带一个Context（父）并返回一个导出的Context（子）和一个CancelFunc（取消函数）。调用CancelFunc函数来取消子和他的孩子，移除父母的参考孩子，并停止任何关联的定时器。调用CancelFunc失败会泄露子和他的孩子直到父取消或者计时器超时。go的审查工具被用在所有的control-flow路径下检查CancelFuncs。

*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				fmt.Println("监控退出，停止了...")
				return
			default:
				fmt.Println("goroutine监控中...")
				time.Sleep(2 * time.Second)
			}
		}
	}(ctx)

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	//为了检测监控过是否停止，如果没有监控输出，就表示停止了
	time.Sleep(5 * time.Second)
}