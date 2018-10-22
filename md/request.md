#`go http`包,模拟客户端`http`请求,处理请求数据与返回数据

## 八种请求方法
   
| 序号        | 方法    |  描述  |
| --------   | -----:   | :----: |
| 1 | GET |请求指定的页面信息，并返回实体主体或者说正文|
| 2 | HEAD | 类似于get请求，只不过返回的响应中没有具体的内容，用于获取报头 |
| 3 | POST | 向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。<br>数据被包含在请求体中。POST请求可能会导致新的资源的建立或已有资源的修改。 |
| 4 | PUT |从客户端向服务器传送的数据取代指定的文档的内容。|
| 5 | DELETE |请求服务器删除指定的页面。|
| 6 | CONNECT |HTTP/1.1协议中预留给能够将连接改为管道方式的代理服务器。|
| 7 | OPTIONS |允许客户端查看服务器的性能。|
| 8 | TRACE|回显服务器收到的请求，主要用于测试或诊断。|


## 请求数据格式 主要有两种

> 测试代码

```go
package main

import (
	"log"
	"net/http"
	"fmt"
)

func main()  {
	db	:=	database{"shoes":50,"socks":5}
	log.Fatal(http.ListenAndServe("127.0.0.1:8000",db))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f",	d)
}

//数据结构
type database map[string]dollars

func (db database) ServeHTTP(w http.ResponseWriter,req *http.Request) {
	//测试代码
}
```

### 1. `queryString` 格式

* 字符串查询模式，简单来说就是链接后面拼接查询字符
* 例如： [https://ZRothschild:123456@studyiris.com:80/detail?id=46&uid=13#jwt](https://ZRothschild:123456@studyiris.com:80/detail?id=46&uid=13#jwt)

> queryString 获取请求值是使用 

```go
	uri,_ := json.Marshal(req.URL)
	w.Write(uri)
	
    req.ParseForm()
    reqForm,_ := json.Marshal(req.Form)
    w.Write(reqForm)
    //{"id":["46"] "uid":["13"]}
```

### 2. `body`  请求体格式

> 请求体格式又可以细分为一下几种

#### 1. `form-data` 格式

* 就是http请求中的multipart/form-data,它会将表单的数据处理为一条消息，以标签为单元，用分隔符分开。
* 既可以上传键值对，也可以上传文件。
* 当上传的字段是文件时，会有Content-Type来说明文件类型；
* content-disposition，用来说明字段的一些信息；
* 由于有boundary隔离，所以multipart/form-data既可以上传文件，也可以上传键值对，
* 它采用了键值对的方式，所以可以上传多个文件。

```go
     //--------------561079950965484784119941 (boundary)
    //Content-Disposition: form-data; name="verifyString"; filename="default.jpg"
    //Content-Type: image/jpeg
    
    
    //----------------------------492053930847914637974767
   // Content-Disposition: form-data; name="lang"
    //
    //go
    //----------------------------492053930847914637974767--
    // Content-Disposition: form-data; name="frame"
    //
    //iris
    //----------------------------150059972171157177176959--
```

> `form-data/multipart/form-data`格式下 `Form`，`PostForm` 将得不到值

```go
    //请求 form-data  键 lang  值 go   键 file 值 default.jpg对象
	req.ParseForm()
	reqForm,_ := json.Marshal(req.Form)
	w.Write(reqForm)
	//{}
	
	postForm,_ := json.Marshal(req.PostForm)
	w.Write(postForm)
	//{}
	
    req.ParseMultipartForm(100)
    multipartForm,_ := json.Marshal(req.MultipartForm)
    w.Write(multipartForm)
    //{
    //    "Value": {
    //        "lang": [
    //            "go"
    //        ]
    //    },
    //    "File": {
    //        "file": [
    //            {
    //                "Filename": "default.jpg",
    //                "Header": {
    //                    "Content-Disposition": [
    //                        "form-data; name=\"verifyString\"; filename=\"default.jpg\""
    //                    ],
    //                    "Content-Type": [
    //                        "image/jpeg"
    //                    ]
    //                },
    //                "Size": 10221
    //            }
    //        ]
    //    }
    //}
    
```

#### 2. `application/x-www-from-urlencoded` 格式

* 会将表单内的数据转换为键值对，比如,lang=go&frame=iris

> `x-www-from-urlencoded` 格式下 `Form`，`PostForm` 可以获取到值，不能传文件格式

```go
    //请求数据 {"lang":"go"}
	req.ParseForm()
	reqForm,_ := json.Marshal(req.Form)
	w.Write(reqForm)
	//{"lang":["go"]}
	postForm,_ := json.Marshal(req.PostForm)
	w.Write(postForm)
	//{"lang":["go"]}
```


#### 3. `raw` 格式

> 可以上传任意格式的文本，可以上传`text`、`json`、`xml`、`html`等

```go
    //请求 raw {"lang":"go"}
    body,_ := ioutil.ReadAll(req.Body)
    defer req.Body.Close()
    w.Write(body)
    //{
    //    "lang": "go"
    //}
```
> `raw`格式下 Form`，`PostForm` 将得不到值

```go
    //请求 raw {"lang":"go"}
	req.ParseForm()
	aa,_ := json.Marshal(req.Form)
	w.Write(aa)
	//{}
	bb,_ := json.Marshal(req.PostForm)
	w.Write(bb)
	//{}
```

#### 4. `binary` 格式

* 相当于Content-Type:application/octet-stream,
* 从字面意思得知，只可以上传二进制数据，
* 通常用来上传文件，由于没有键值，所以，一次只能上传一个文件。

```go
    // 传二进制
	req.ParseForm()
	reqForm,_ := json.Marshal(req.Form)
	w.Write(reqForm)
	//{}

	postForm,_ := json.Marshal(req.PostForm)
	w.Write(postForm)
	//{}

	req.ParseMultipartForm(100)

	multipartForm,_ := json.Marshal(req.MultipartForm)
	w.Write(multipartForm)
	// null
```

### 3. 获取请求数据分析总结

#### `queryString` 格式

> 不区分请求方法，区分请求格式

* 以`queryString`形式请求数据，可以用`req.URL` 或者 `req.Form`切记要先调用 `req.ParseForm()`

#### `body` 格式

> `PostForm`例外只适用`POST`请求

* `form-data/multipart/form-data`格式下`req.Form`，`PostForm`将返回`{}`，可以用`req.MultipartForm`
* `x-www-from-urlencoded`可用`req.Form`，`PostForm`
* `raw`不支持`req.Form`，`PostForm`
* `binary`不支持`req.Form`，`PostForm`

> **如果你不想费事劳神，接收读取请求数据统一使用req.Body,他可以接收任何数据格式**

## 数据处理

### 请求数据处理

* `req.Form`，`PostForm`返回值为`url.Values(map[string][]string)`类型,有`Get`,`Set`,`Add`,`Del`方法
* `PostForm`返回值为`*multipart.Form`类型,结构如下，详细请参考`multipart`包

```go
type Form struct {
	Value map[string][]string
	File  map[string][]*FileHeader
}
```

* `req.URL`返回值为`*url.URL`类型,结构如下,可以参考我上一篇文章，详细请参考 `url`包

```go
//[scheme:][//[userinfo@]host][/]path[?query][#fragment]
type URL struct {
	Scheme     string
	Opaque     string    // 编码的不透明数据
	User       *Userinfo // 用户名和密码信息
	Host       string    // 主机或主机：端口
	Path       string    // 路径（相对路径可能省略前导斜杠）
	RawPath    string    // 编码路径提示（请参阅EscapedPath方法）
	ForceQuery bool      // 即使RawQuery为空，也附加一个查询（'？'）
	RawQuery   string    // 编码的查询值，没有'？'
	Fragment   string    // 参考片段，'＃'后面的值
}
```
* `req.Body`返回值为`io.ReadCloser`接口，结构如下,详细请参考 `io`包

```go
type ReadCloser interface {
	Reader
	Closer
}
```

> 一般的处理流程为

```go
//返回byte 类型，然后在做转换
body,_ := ioutil.ReadAll(req.Body)
defer req.Body.Close()
```

### 返回数据处理

> 可以自定义返回数据格式

```go
// 要输出json记得加头信息?
w.Header().Add("Content-Type", "application/json") 
```

* 数据组装然后通过 `fmt.Fprintf(w,"%s:\n",str)`返回
* 数据组装然后通过 `w.Write(test)`返回


```go
	body,_ := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	w.Write(body)
	//str  为输出字符串
	fmt.Fprintf(w,"%s:\n",str)
```

