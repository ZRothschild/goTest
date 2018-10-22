## `go`语言`net/http`包 http server 服务

> 个人觉得Go圣经写的非常出色，所有在这里抄了一波。建议初学者先看几遍Go圣经。

### `http.Handler`接口 

#### 1.简单实现`http.Handler`接口的服务器**API**

```go
package	http

type Handler interface {
	ServeHTTP(w ResponseWriter,r *Request)
}

func ListenAndServe(address string,h Handler) error
```

* **ListenAndServe**函数需要一个例如**localhost:8000**的服务器地址,
* 和一个所有请求都可以分派的**Handler**接口实例。**(注意是所有请求都可以分派)**
* 它会一直运行,直到这个服务因为一个错误而失败(或者启动失败),它的返回值一定是一个非空的错误。
* 想象一个电子商务网站,为了销售它的数据库将它物品的价格映射成美元。
* 下面这个程序可能是能想到的最简单的实现了。
* 它将库存清单模型化为一个命名为database的map类型,我们给这个类型一个`ServeHttp`方法,这样它可以满足`http.Handler`接口。
* 这个handler会遍历整个map并输出物品信息。

> 文件 `main.go`

```go
package main

import (
	"log"
	"net/http"
	"fmt"
)

func main()  {
	db := database{"shoes":50,"socks":5}
	log.Fatal(http.ListenAndServe(":8000",db))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f",	d)
}

//数据结构
type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter,req *http.Request) {
	for item,price := range db {
		fmt.Fprintf(w,"%s: %s\n",item,price)
	}
}
```
>  执行 `main.go` 浏览器访问，`localhost:8000` 会输出

```go
//shoes: %!s(main.dollars=50)
//socks: %!s(main.dollars=5)
```

> 当然也可以用`http`方法请求，例如用**postman**`POST`返回结果同上

#### 2.服务器定义多个不同的URL

* 第一部分的服务器不考虑URL只能为每个请求列出它全部的库存清单。
* （回忆上面一句话）一个所有请求都可以分派的`Handler`接口实例
* 更真实的服务器会定义多个不同的URL,每一个都会触发一个不同的行为。
* 让我们使用/list来调用已经存在的这个行为并且增加另一个`/price`调用表明单个货品的价格
* 或像这样`/price?item=socks`来指定一个请求参数。

```go
package main

import (
	"log"
	"net/http"
	"fmt"
)

func main()  {
	db := database{"shoes":50,"socks":5}
	
	log.Fatal(http.ListenAndServe(":8000",db))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f",	d)
}

//数据结构
type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter,req *http.Request) {
	switch req.URL.Path {
	case "/list":
		for item,price := range db {
			fmt.Fprintf(w,"%s: %s\n",item,price)
		}
	case "/price":
		item := req.URL.Query().Get("item")
		price,ok :=	db[item]
		if !ok {
			w.WriteHeader(http.StatusNotFound)//404
			fmt.Fprintf(w,"no such item:%q\n",item)
			return
		}
		fmt.Fprintf(w,"%s\n",price)
	default:
		w.WriteHeader(http.StatusNotFound)//404
		fmt.Fprintf(w,"no such item:%q\n",req.URL)
	}
}
```

* 现在`handler`基于`URL`的路径部分`(req.URL.Path)`来决定执行什么逻辑。
* 如果这个`handle`不能识别这个路径,它会通过调用`w.WriteHeader(http.StatusNotFound)`返回客户端一个HTTP错误;
* 这个检查应该在向`w`写入任何值前完成。`w.WriteHeader`在另外一篇文章有详解
* `http`包有一个错误处理方法http.Error函数:

```go
    msg := fmt.Sprintf("no such page: %s\n",req.URL)
    http.Error(w,msg,http.StatusNotFound) //404
```

* `/price`的case会调用URL的Query方法来将HTTP请求参数解析为一个map,
* 或者更准确地说一个`net/url`包中`url.Values()`类型的多重映射。
* 然后找到第一个item参数并查找它的价格。如果这个货品没有找到会返回一个错误。

> 上面已经举过例子了，我相信大家不言自明

#### 3.`web`服务器请求多路器ServeMux

* 对于现在的`web`站点来说，请求路径是10的n次方，很明显一直`case`下去会要老命的
* 所有东西都往里面方，代码与业务也不好管理。
* 所以我们需要一个服务器请求多路器`ServeMux`

```go
package main

import (
	"log"
	"net/http"
	"fmt"
)
func main() {
	db := database{"shoes":50,"socks":5}
	mux := http.NewServeMux()
	mux.Handle("/list",http.HandlerFunc(db.list))
	mux.Handle("/list",http.HandlerFunc(db.list))
	log.Fatal(http.ListenAndServe(":8000",mux))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f",	d)
}
//数据结构
type database map[string]dollars
//list
func (db database) list(w http.ResponseWriter, req *http.Request){
	for item,price := range db {
		fmt.Fprintf(w,"%s: %s\n",item,price)
	}
}
//price
func (db database) price(w http.ResponseWriter, req *http.Request){
	item := req.URL.Query().Get("item")
	price,ok :=	db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)//404
		fmt.Fprintf(w,"no such item:%q\n",item)
		return
	}
	fmt.Fprintf(w,"%s\n",price)
}
```
* 让我们关注这两个注册到handlers上的调用。
* 第一个db.list是一个方法值,它是下面这个类型的值
> 如果您不了解方法值，我建议您可以看**Go圣经**,关于函数，方法那两段段，前提你要有起码的基础

```go
func (w http.ResponseWriter, req *http.Request)
```

* 也就是说`db.list`的调用会援引一个接收者是`db`的`database.list`方法。
* 虽然`db.list`是一个实现了`handler`类似行为的函数,但是因为它没有方法
* 因为它不满足`http.Handler`接口并且不能直接传给`mux.Handle`。
* 语句`http.HandlerFunc(db.list)`是一个转换而非一个函数调用,因为`http.HandlerFunc`是一个类型。
> 多说一句类型转换，类型断言这是基础哦， 它有如下的定义:

```go
package	http

type HandlerFunc (w ResponseWriter,r *Request)

func (f HandlerFunc) ServeHTTP(w http.ResponseWriter,req *http.Request) {
	f(w,r)
}
```

* `HandlerFunc`显示了在`Go`语言接口机制中一些不同寻常的特点。
* 这是一个有实现了接口`http.Handler`方法的函数类型。
* `ServeHTTP`方法的行为调用了它本身的函数。
* 因此`HandlerFunc`是一个让函数值满足一个接口的适配器,这里函数和这个接口仅有的方法有相同的函数签名。
* 实际上,这个技巧让一个单一的类型例如`database`以多种方式满足`http.Handler`接口:
* 一种通过它的`list`方法,一种通过它的`price`方法等等。
* 因为handler通过这种方式注册非常普遍,ServeMux有一个方便的HandleFunc方法,
* 它帮我们简化handler注册代码成这样:

```go
mux.HandleFunc("/list", db.list)
mux.HandleFunc("/price", db.price)
```

* 从上面的代码很容易看出应该怎么构建一个程序,它有两个不同的`web`服务器监听不同的端口的,
* 并且定义不同的`URL`将它们指派到不同的`handler`。
* 我们只要构建另外一个`ServeMux`并且在调用一次`ListenAndServe`(可能并行的)。
* 但是在大多数程序中,一个`web`服务器就足够了。
* 此外,在一个应用程序的多个文件中定义`HTTPhandler`也是非常典型的,
* 如果它们必须全部都显示的注册到这个应用的`ServeMux`实例上会比较麻烦。
* `net/http`包提供了一个全局的`ServeMux`实例`DefaultServerMux`和包级别的`http.Handle`和`http.HandleFunc`函数。
* 现在,为了使用`DefaultServeMux`作为服务器的主`handler`,我们不需要将它传给`ListenAndServe`函数;`nil`值就可以工作。

```go
func main() {
    db := database{"shoes":50,"socks":5}
    http.HandleFunc("/list",db.list)
    http.HandleFunc("/price",db.price)
    log.Fatal(http.ListenAndServe("localhost:8000",	nil))
}
```

> 简单的服务完成，当然`http`服务并没有那么简单，这只是入门。
