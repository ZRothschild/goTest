# go context 包

### 介绍

>  在`Go`服务器中，每个传入的请求都在其自己的`goroutine`中处理。
>  请求处理程序通常会启动其他`goroutine`来访问后端，例如数据库和RPC服务
> 处理请求的`goroutine`集合通常需要访问特定于请求的值，例如最终用户的身份，授权令牌和请求的截止日期。
> 当请求被取消或超时时，处理该请求的所有`goroutine`都应该快速退出，以便系统可以回收它们正在使用的任何资源。
> 在`Google`，我们开发了一个[context](https://golang.org/pkg/context)包，可以轻松地将请求范围的值，取消信号和`API`边界的截止日期传递给处理请求所涉及的所有`goroutine`。
> 该软件包作为上下文公开提供。 本文介绍如何使用该包并提供完整的工作示例。**注意go1.7版本开始有这个包**

### context解释

1. 包`context`定义了`Context`类型,并包含`deadlines`(结束时间),`cancelation signals`(取消信号),和其他的请求`api`传送的值

2. 对服务器的请求应该创建一个`Context`,服务器发出的外向请求应该接受`Context`.链式的函数调用之间必须传递`Context`,随意的更换它使用一个导出
`Context`使用`WithCancel`、`WithDeadline`、`WithTimeout`、`WithValue`。当一个`Context`取消所有从它衍生的`Contexts`都会被取消.

3. `WithCancel`、`WithDeadline`、`WithTimeout`函数携带一个Context(父)并返回一个衍生的Context(子)和一个`CancelFunc`(取消函数).
调用CancelFunc函数来取消子和他的孩子,移除父母的参考孩子,并停止任何关联的定时器。调用`CancelFunc`失败会泄露子和他的孩子直到父取消或者计时器超时.
`go`的审查工具被用在所有的`control-flow`路径下检查`CancelFuncs`.

### 使用上下文的程序应遵循这些规则来保留接口跨包一致，并启用静态分析工具来检查上下文：

1. 不要将上下文存储在结构类型中; 相反,传递一个上下文,明确地显示需要它的每个函数.并且应该作为第一个参数,通常命名为ctx：

```go
    func DoSomething(ctx context.Context, arg Arg) error {
        // ... use ctx ...
    }
```

2. 即使函数允许.也不要传递`nil Context`.传递`context.TODO`如果您不确定要使用哪个上下文.

3. 仅将上下文值用于转换进程的请求范围数据API,不用于将可选参数传递给函数.

4. 可以将相同的`Context`传递给在不同`goroutine`中运行的函数;上下文可以安全地同时使用多个`goroutine`.

### Context 结构体

```go
    // A Context carries a deadline, cancelation signal, and request-scoped values
    // across API boundaries. Its methods are safe for simultaneous use by multiple
    // goroutines.
    type Context interface {
        // Done returns a channel that is closed when this Context is canceled
        // or times out.
        Done() <-chan struct{}
    
        // Err indicates why this context was canceled, after the Done channel
        // is closed.
        Err() error
    
        // Deadline returns the time when this Context will be canceled, if any.
        Deadline() (deadline time.Time, ok bool)
    
        // Value returns the value associated with key or nil if none.
        Value(key interface{}) interface{}
    }
```

- Done()，Done方法返回一个channel，作为代表Context运行的函数的取消信号(执行cancel函数)：当channel关闭时，函数应该放弃它们的工作并返回。 Err方法返回一个错误(Context被取消的原因)
- Err()，返回一个错误。该context为什么被取消掉。
- Deadline()，返回截止时间和ok，Deadline方法允许函数确定它们是否应该开始工作;如果剩下的时间太少，可能就不值得了。代码也可以使用截止时间来设置`I/O`操作的超时。。
- Value()，返回值，Value允许Context携带请求范围的数据。该数据必须是安全的，以便多个goroutine同时使用。

### Context 所以方法

```go
    func Background() Context
    func TODO() Context
    
    func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
    func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
    func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
    func WithValue(parent Context, key, val interface{}) Context
```

### 例子

- 例一

```go
package main

import (
	"time"
	"context"
	"fmt"
)

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
	time.Sleep(5 * time.Second)
}
```

- 例二

```go
package main

import (
	"time"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	go watch(ctx,"【监控1】")
	go watch(ctx,"【监控2】")
	go watch(ctx,"【监控3】")

	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
	fmt.Println(1111)
}

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Println(name,"监控退出，停止了...")
			f,_ := os.Create(name+".txt")
			f.Write([]byte(name))
			f.Close()
			return
		default:
			fmt.Println(name,"goroutine监控中...")
			time.Sleep(2 * time.Second)
		}
	}
}
```

- 例三伴随传递数据

```go
package main

import (
	"time"
	"context"
	"fmt"
)

var key string="name"
var keys string="names"

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	//附加值
	valueCtx:=context.WithValue(ctx,key,"【监控1】")
	go watch(valueCtx)
	aa := context.WithValue(valueCtx,keys,"【监控3】")
	go watch(aa)
	time.Sleep(10 * time.Second)
	fmt.Println("可以了，通知监控停止")
	cancel()
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
```

-  例四伴随传递数据，使用`WithDeadline`

```go
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
```

### 总结

`context`包是用来管理进程以及衍生出的协程。打一个比方，一个项目三组人，项目相当于进程，每一组相对于一个线程，组里面每个人相当于协程。
当一个组任务完成了，组员们可以停下你们的工作了。那么现实中就是需要，组长告知大家，所有人停止任务。同样协程也需要，一个信号告诉他们停止任务。
`context`包就起到通知的功能。项目进程持有一个顶级`A Context`，每个组持有`A Context`衍生的`B,C,D Context`,组成员持有组衍生`X,Y,Z Context`.