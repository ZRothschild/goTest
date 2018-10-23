# go context 包

### context解释

1. 包`context`定义了`Context`类型,并包含`deadlines`(结束时间),`cancelation signals`(取消信号),和其他的请求`api`范围的值,进程

2. 对服务器的请求应该创建一个`Context`,服务器发出的外向请求应该接受`Context`.链式的函数调用之间必须传递`Context`,随意的更换它使用一个导出
`Context`使用`WithCancel`、`WithDeadline`、`WithTimeout`、`WithValue`。当一个`Context`取消所有从它导出的`Contexts`都会被取消.

3. `WithCancel`、`WithDeadline`、`WithTimeout`函数携带一个Context(父)并返回一个导出的Context(子)和一个`CancelFunc`(取消函数).
调用CancelFunc函数来取消子和他的孩子,移除父母的参考孩子,并停止任何关联的定时器。调用`CancelFunc`失败会泄露子和他的孩子直到父取消或者计时器超时.
`go`的审查工具被用在所有的`control-flow`路径下检查`CancelFuncs`.

### 使用上下文的程序应遵循这些规则来保留接口跨包一致，并启用静态分析工具来检查上下文 传播：

1. 不要将上下文存储在结构类型中; 相反,传递一个上下文,明确地显示需要它的每个函数.语境应该是第一个参数,通常命名为ctx：

```
    func DoSomething(ctx context.Context, arg Arg) error {
        // ... use ctx ...
    }
```

2. 即使函数允许.也不要传递`nil Context`.传递`context.TODO`如果您不确定要使用哪个上下文.

3. 仅将上下文值用于转换进程的请求范围数据API,不用于将可选参数传递给函数.

4. 可以将相同的`Context`传递给在不同`goroutine`中运行的函数;上下文可以安全地同时使用多个`goroutine`.