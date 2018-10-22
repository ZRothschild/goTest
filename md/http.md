## `go http`包 `http`请求
## 响应体 `w http.ResponseWriter`
#### `w http.ResponseWriter` 接口类型 数据结构
```go
    //HTTP处理程序使用ResponseWriter接口构造一个HTTP响应。
    //在Handler.ServeHTTP方法返回之后，可能无法使用ResponseWriter
type ResponseWriter interface {
    //Header返回的header map将被 WriteHeader 返回。 
    //Header map 也是其中的机制处理程序可以设置HTTP trailers。
    
    //调用WriteHeader或（Write）后改变header map是无效的，
    //除非重定义的标题 headers 也就是 headers 
    
    // There are two ways to set Trailers. The preferred way is to
    // predeclare in the headers which trailers you will later  // send by setting the "Trailer" header to the names of the
    // trailer keys which will come later. In this case, those
    // keys of the Header map are treated as if they were
    // trailers. See the example. The second way, for trailer
    // keys not known to the Handler until after the first Write,
    // is to prefix the Header map keys with the TrailerPrefix
    // constant value. See TrailerPrefix.
    //要禁止隐式响应标头（例如“Date”），请设置他们的价值为零。
	Header() Header

	//Write将数据作为HTTP回复的一部分写入连接。
	//如果尚未调用WriteHeader，Write调用在写入数据之前写入WriteHeader http.StatusOK
    //如果是Header不包含Content-Type行，
    //Write添加Content-Type集将最初的512字节写入数据传递给的结果 DetectContentType
	//根据HTTP协议版本和客户端，调用Write或WriteHeader可能会阻止将来读取Request.Body
	//对于HTTP/1.x请求，处理程序应该读取任何内容,在编写响应之前需要请求正文数据。 
    //一旦headers已被刷新（由于显式的Flusher.Flush调用或写入足够的数据来触发刷新），
    //即请求体可能无法使用 对于HTTP/ 2请求，Go HTTP服务器允许,处理程序在同时继续读取请求体写回复 
    //但是，可能不支持此类行为由所有HTTP / 2客户端。处理者应在写作前阅读可以最大化兼容性。
	Write([]byte) (int, error)
	
	//WriteHeader发送带有状态代码的HTTP响应头。
    //如果未显式调用WriteHeader，则首先调用Write
    //将触发隐式WriteHeader（http.StatusOK）。
    //因此显式调用WriteHeader主要用于发送错误代码
	WriteHeader(int)
}
```
### `w http.ResponseWriter`接口详解
#### `Header()` 方法
`Header()`方法返回`Header`类型结构如下
```go
//Header 表示HTTP Header中的键值对。
type Header map[string][]string
```
##### `Header`有如下方法七个方法
###### 1.`Add`添加
```go
//Add将键值对添加到 Header map中。
//它附加到与key关联的任何现有值。 注意 是追加 而不是覆盖
func (h Header) Add(key, value string) {
	textproto.MIMEHeader(h).Add(key, value)
}
```
###### 2.`Set`设置
```go
//Set设置与key关联的标题条目单个元素值。 
//它取代了任何现有的与key关联的值。注意 存在即覆盖
func (h Header) Set(key, value string) {
	textproto.MIMEHeader(h).Set(key, value)
}
```
###### 3.`Get`获取
```go
//获取与给定键关联的第一个值。不区分大小写; 
//使用textproto.CanonicalMIMEHeaderKey,规范化提供的key。
//如果没有与键关联的值，则Get返回空字符串。要访问key多个值，或使用非规范key，直接访问map。
//注意 如果想得到key对用的所有值，可以直接像map一样访问
func (h Header) Get(key string) string {
	return textproto.MIMEHeader(h).Get(key)
}
```
###### 4.`get`获取
```go
// get就像Get，但key必须已经是CanonicalHeaderKey形式。
func (h Header) get(key string) string {
	if v := h[key]; len(v) > 0 {
		return v[0]
	}
	return ""
}
```
###### 5.`Del`删除
```go
// Del删除与key关联的值。
func (h Header) Del(key string) {
	textproto.MIMEHeader(h).Del(key)
}
```
###### 6.`Write`以特定格式写入
```go
// Write以有线格式写入header。
func (h Header) Write(w io.Writer) error {
	return h.WriteSubset(w, nil)
}
```
###### 7.clone`复制一份header
```go
func (h Header) clone() Header {
	h2 := make(Header, len(h))
	for k, vv := range h {
		vv2 := make([]string, len(vv))
		copy(vv2, vv)
		h2[k] = vv2
	}
	return h2
}
```
> 这七个方法就是对`Header`的基本操作 示例如：
```go
//自定义头部信息
w.Header().Add("name", "ZRothschild")
//数据格式头json
w.Header().Add("Content-Type", "application/json")
//json 格式数据
test,_ := json.Marshal(map[string]string{"title":"标题"})
//返回json 格式数据
w.Write(test)
// 要输出HTML记得加头信息
w.Header().Add("Content-Type", "text/html") 
//html/text
w.Write([]byte("<h1>Hello World</h1>"))
```
####  `Write([]byte) (int, error)` 方法
`Write`方法将会把其参数写入返回正文，`int`返回类型表示写入的字节数
```go
//将会返回我爱golang intNum 写入字节数
intNum,err := w.Write([]byte("我爱golang"))
```
####  `WriteHeader(int)` 方法
`WriteHeader`用于设置返回状态，在`Write`方法前设置，不设置默认`http.StatusOK`
如果你希望请求返回301表示重定向那么你应该这样做
```go
jsonData,_ := json.Marshal(map[string]string{"msg":"接收后请重定向"})
w.WriteHeader(301)
w.Write(test)
```
## 请求体 `req *http.Request`
####  `req *http.Request`结构体类型 数据结构
```go
type Request struct {
	// Method指定HTTP方法（GET，POST，PUT等）。
    //对于客户端请求，空字符串表示GET。
	Method string
	
	//URL指定要求的URI（对于服务器请求）或访问的URL（用于客户端请求）
	//请求服务器与访问客户端的资源地址。
	//对于服务器请求，URL是从URI解析的出来的，由RequestURI中提供，因为数据存储在RequestURI中。
	//对于大多数请求，Path和RawQuery以外的字段都是空 （参见RFC 2616，第5.1.2节）
	//对于客户端请求，URL的host指定连接服务器，要指定http请求header头，host字段的名称。
	//header["host"] = host
	URL *url.URL
	
	//传入服务器请求的协议版本。对于客户端请求，将忽略这些字段。
    //HTTP 客户端代码始终使用HTTP / 1.1或HTTP / 2。
    //有关详细信息，请参阅传输文档。
	Proto      string // "HTTP/1.0"
	ProtoMajor int    // 1
	ProtoMinor int    // 0

	//Header包含收到的请求Header字段
    //由服务器或由客户端发送。
	//如果服务器收到Header的请求，header 用法
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"Bar", "two"},
	//	}

	//对于传入请求，Host header被提升为Request.Host字段并从Header map中删除。
	//HTTP定义标题名称不区分大小写。
	//该请求解析器通过使用CanonicalHeaderKey实现这一点，
    //第一个字符和后面的任何字符连字符首字母大写，其余小写。
	//对于客户端请求，某些headers，
	//例如Content-Length并在需要时自动写入ConnectionHeader中的值可能会被忽略。
	//请参阅文档
    //用于Request.Write方法。
	Header Header
	
	//正文是请求的正文。
    //对于客户端请求，nil表示请求没有正文，例如GET请求。
    //HTTP客户端的传输负责调用Close方法。
    //对于服务器请求，Request Body始终为非零
    //但是当没有正文时会立即返回EOF。
    //服务器将关闭请求正文。ServeHTTP处理程序不需要关闭。
	Body io.ReadCloser
	
	//GetBody定义了一个可选的func来返回一个新副本body。
	//当重定向需要时，它用于客户端请求多次读取body。 
	//仍然使用GetBody需要设置Body。
    //对于服务器请求，它是未使用的。
	GetBody func() (io.ReadCloser, error)
	// ContentLength记录关联内容的长度。
    //值-1表示长度未知。
    //值> = 0表示给定的字节数可以从Body中读取。
    //对于客户端请求，带有非零主体的值为0也被视为未知。
	ContentLength int64
	
	//TransferEncoding列出从最外层到的传输编码最里面 空列表表示“身份”编码。
    //通常可以忽略TransferEncoding;分块编码是在发送时自动添加和删除接收请求
	TransferEncoding []string
	
	//关闭表示是否在之后关闭连接
    //回复此请求（对于服务器）或发送此请求后 请求并阅读其响应（针对客户）。
	//对于服务器请求，HTTP服务器会自动处理此请求处理程序不需要此字段。
	//对于客户端请求，设置此字段可防止重复使用对同一主机的请求之间的TCP连接，
	// 如同设置了Transport.DisableKeepAlives。
	Close bool
	
	//对于服务器请求Host指定其上的主机寻找网址。 
	//根据RFC 2616，header host键的值或URL本身中给出的主机名。
    //它可能是“host：port”的形式。对于国际域名名称，主机可以是Punycode或Unicode形式。
    //使用 golang.org/x/net/idna将其转换为任何格式如果需要
	//对于客户端请求主机可选择覆盖主机 要发送的标头。 
	//如果为空，则Request.Write方法使用URL.Host的值。 
	//Host 可能包含国际域名。
	Host string
	
	//表单包含已解析的表单数据，包括URL字段的查询参数和POST或PUT表单数据。
    //此字段仅在调用ParseForm后可用。
    //HTTP客户端忽略Form并改为使用Body。
	Form url.Values

	//PostForm包含来自POST，PATCH的已解析表单数据，或PUT体参数。
    //此字段仅在调用ParseForm后可用。
    //HTTP客户端忽略PostForm并使用Body。
	PostForm url.Values
	
	//MultipartForm是解析更广扩展部分表单，包括文件上载。
    //此字段仅在调用ParseMultipartForm后可用。
    //HTTP客户端忽略MultipartForm并改为使用Body。
	MultipartForm *multipart.Form
	
	//Trailer指定在请求体后发送的其headers。
	//对于服务器请求，Trailer map 最初只包含 Trailer密钥，nil值。
	//（客户声明哪些Trailer稍后会发送。）当处理程序从Body读取时，它必须 不是参考Trailer。
	//读取Body后返回EOF，Trailer可以由客户端再次读取并包含非零值（如果已发送）。
	//对于客户端请求，必须将Trailer初始化为map包含的以后发送的Trailer键。 
	//值可以是零或最终值。 ContentLength必须为0或-1，以发送分块请求。
    //发送HTTP请求后，可以更新map值 读取请求正文。 
    //一旦body返回EOF，调用者则者必须 不要改变Trailer。
	//很少有HTTP客户端，服务器或代理支持HTTP Trailer。
	Trailer Header
	
	//RemoteAddr允许HTTP服务器和其他软件进行记录发送请求的网络地址，通常用于记录。
	//ReadRequest并未填写此字段没有定义的格式。 
	//此包中的HTTP服务器在调用之前将RemoteAddr设置为“IP：port”地址处理程序。
	//HTTP客户端忽略此字段。
	RemoteAddr string
	
	//RequestURI是未经修改的Request-URI
    //客户端发送的Request-Line（RFC 2616，第5.1节）
    //到服务器 通常应该使用URL字段。
    //在HTTP客户端请求中设置此字段是错误的。
	RequestURI string
	
    //TLS允许HTTP服务器和其他软件进行记录有关请求的TLS连接的信息。
    //ReadRequest未填写此字段。
    //此程序包中的HTTP服务器为其设置字段
    //在调用处理程序之前启用TLS连接;否则它会使字段为零。
    // HTTP客户端忽略此字段。
	TLS *tls.ConnectionState
	
	//Cancel是一个可选通道，其闭包表示客户端请求应被视为已取消。 
	//不是所有的实现 RoundTripper 支持取消。
    //对于服务器请求，此字段不适用。
    //不推荐使用：使用Context和WithContext方法来替代。 
    //未定义将不会有作用，是否遵守取消。
	Cancel <-chan struct{}
	
	//此字段仅在客户端重定向期间填充。
	Response *Response
	
	//ctx是客户端或服务器上下文。它应该只通过使用WithContext复制整个请求来修改。
    //为了防止人们使用Context应为错误是不可预知的,改变同一请求的调用者持有的上下文。
	ctx context.Context
}
```
### `req *http.Request` 详解
####  `Method string` `HTTP`请求方法
`Method` 大致有八种
| 序号        | 方法    |  描述  |
| -------- | -----:   | :----: |
| 1 | GET |请求指定的页面信息，并返回实体主体或者说正文|
| 2 | HEAD | 类似于get请求，只不过返回的响应中没有具体的内容，用于获取报头 |
| 3 | POST | 向指定资源提交数据进行处理请求（例如提交表单或者上传文件）。<br>数据被包含在请求体中。POST请求可能会导致新的资源的建立或已有资源的修改。 |
| 4 | PUT |从客户端向服务器传送的数据取代指定的文档的内容。|
| 5 | DELETE |请求服务器删除指定的页面。|
| 6 | CONNECT |HTTP/1.1协议中预留给能够将连接改为管道方式的代理服务器。|
| 7 | OPTIONS |允许客户端查看服务器的性能。|
| 8 | TRACE|回显服务器收到的请求，主要用于测试或诊断。|
当前`HTTP`请求方法，当以`GET`请求即返回`GET`，当以`POST`请求即返回`POST`
#### 请求资源地址`req.URL *url.URL`
##### `req.URL`数据结构
> 注意`req.URL`与`queryString`关联
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
```go
package main

import (
	"fmt"
	"net/url"
	"strings"
)

func main() {
	//我们将解析这个 URL
	s := "https://ZRothschild:123456@studyiris.com:80/detail?id=46&uid=13#jwt"
	urlParse, _ := url.Parse(s)
	//Scheme协议格式
	fmt.Printf("Scheme %s\n",urlParse.Scheme)
	//输出 https

	//Opaque是否透明
	fmt.Printf("Opaque %s\n",urlParse.Opaque)
	//输出空

	//User包含了所有的认证信息，这里调用 Username和Password方法来获取独立值。
	fmt.Printf("User %s\n",urlParse.User)
	//输出 ZRothschild:123456

	//用户名
	fmt.Printf("User Username %s\n",urlParse.User.Username())
	//输出 ZRothschild

	//密码
	pass,_:= urlParse.User.Password()
	fmt.Printf("User Password %s\n",pass)
	//输出 123456

	//Host同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split()
	fmt.Printf("Host %s\n",urlParse.Host)
	//输出 studyiris.com:80
	host := strings.Split(urlParse.Host, ":")
	fmt.Printf("Host %s\n",host[0])
	//输出 studyiris.com
	fmt.Printf("Port %s\n",host[1])
	//输出 80

	//路径
	fmt.Printf("Path %s\n",urlParse.Path)
	//输出 /detail

	fmt.Printf("RawPath %s\n",urlParse.RawPath)
	//输出 为空
	//当url 加密以后 如 https://ZRothschild:123456@studyiris.com:80/article%2fdetail?id=46&uid=13#jwt
	//则会输出  /article%2fdetail 他与Opaque有关
	
	fmt.Printf("ForceQuery %v\n",urlParse.ForceQuery)
	//输出 false
	//当 s = "https://ZRothschild:123456@studyiris.com:80/detail?"
	//输出 true

	fmt.Printf("RawQuery %s\n",urlParse.RawQuery)
	//输出 id=46&uid=13
	//如果想得到键值对可以
	rawMap, _ := url.ParseQuery(urlParse.RawQuery)
	fmt.Printf("rawMap %v\n",rawMap)
	//输出 map[id:[46] uid:[13]]
	//rawMap["id"][0] === 46  rawMap["uid"][0] === 13

	fmt.Printf("Fragment %s\n",urlParse.Fragment)
	//输出 jwt
}
```
**假设请求 `https://ZRothschild:123456@studyiris.com:80/detail?id=46&uid=13#jwt`**
###### 1.`Scheme`协议格式   
    * http 协议格式
    * https 协议格式
    * ftp 协议格式 
```go
    scheme := req.URL.Scheme
    fmt.Printf("Scheme %s\n",scheme)
    //输出  https 
```
###### 2.`Opaque`是否加透明加密
```go
    opaque := req.URL.Opaque
    fmt.Printf("Opaque %s\n",opaque)
    //输出  空
```
###### 3.`User`登录密码用户名
```go
    user := req.URL.User
    fmt.Printf("User %v\n",user)
    //输出  ZRothschild:123456

	fmt.Printf("User Username %s\n",urlParse.User.Username())
	//输出 ZRothschild

	//密码
	pass,_:= urlParse.User.Password()
	fmt.Printf("User Password %s\n",pass)
	//输出 123456
```
###### 4.`Host`主机或主机：端口
```go
    hostStr := req.URL.Host
    //Host同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split()
    fmt.Printf("Host %s\n",hostStr)
    //输出 studyiris.com:80
    host := strings.Split(host, ":")
    fmt.Printf("Host %s\n",host[0])
    //输出 studyiris.com
    fmt.Printf("Port %s\n",host[1])
    //输出 80
```
###### 5.`Path`请求路径`/detail`
```go
	fmt.Printf("Path %s\n",urlParse.Path)
	//输出 /detail  
```

###### 6. `RawPath`
```go
	fmt.Printf("RawPath %s\n",urlParse.RawPath)
	//输出
```
###### 7.`ForceQuery` 
```go
	fmt.Printf("ForceQuery %v\n",urlParse.ForceQuery)
	//输出 false
	//当 s = "https://ZRothschild:123456@studyiris.com:80/detail?"
	//输出 true
```
###### 8.`RawQuery`问号以后的值不包括问号
```go
	fmt.Printf("RawQuery %s\n",urlParse.RawQuery)
	//输出 id=46&uid=13
	//如果想得到键值对可以
	rawMap, _ := url.ParseQuery(urlParse.RawQuery)
	fmt.Printf("rawMap %v\n",rawMap)
	//输出 map[id:[46] uid:[13]]
	//rawMap["id"][0] === 46  rawMap["uid"][0] === 13
```
###### 9.`ForceQuery` 
```go
	fmt.Printf("Fragment %s\n",urlParse.Fragment)
	//输出 jwt
```
> 提示`req.URL`还是实现论很多方法 详情参见`url`包
#### `Header`头`req.Header`Header类型（`map[string][]string`）
返回`map[string][]string`数据
```go
    fmt.Printf("Header %s\n",req.Header)
    //输出 Header := map[
    //  Connection:[keep-alive] 
    //  User-Agent:[Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 Chrome/65.0.3325.146 Safari/537.36]
    //  Accept:[image/webp,image/apng,image/*,*/*;q=0.8]
    //  Referer:[http://localhost:8000/test/aa?test=11&aa=112]
    //  Accept-Encoding:[gzip, deflate, br]
    //  Accept-Language:[zh-CN,zh;q=0.9]
    // ]
```
#### 请求体`Body io.ReadCloser`
* `body`是请求的正文。对于客户端请求，`nil`表示请求没有正文，例如`GET`请求。
* `HTTP`客户端的传输负责调用`Close`方法。
> **注意**//HTTP客户端忽略`Form`，`PostForm`，`MultipartForm`并改为使用Body
```go
	body,_ := ioutil.ReadAll(req.Body)
	defer defer req.Body.Close()
	w.Write(body)
```
* 对于服务器请求，`Request Body`始终为非零,但是当没有正文时会立即返回`EOF`。
* 服务器将关闭请求正文。`ServeHTTP`处理程序不需要关闭。
#### 复制拷贝 Body `GetBody func() (io.ReadCloser, error)`
* `GetBody`定义了一个可选的`func`来返回一个新副本`body`。
* 当重定向需要时，它用于客户端请求多次读取`body`。 
* 仍然使用`GetBody`需要设置`Body`。对于服务器请求，它是未使用的。
#### 关联内容的长度`ContentLength int64`
* `ContentLength`记录关联内容的长度。值-1表示长度未知。
* 值> = 0表示给定的字节数可以从`Body`中读取。对于客户端请求，带有非零主体的值为0也被视为未知。
#### `TransferEncoding []string`
* `TransferEncoding`列出从最外层到的传输编码最里面 空列表表示“身份”编码。
* 通常可以忽略`TransferEncoding`;分块编码是在发送时自动添加和删除接收请求
#### 关闭请求 `Close bool`
* 关闭表示是否在之后关闭连接,回复此请求（对于服务器）或发送此请求后 请求并阅读其响应（针对客户）。
* 对于服务器请求，`HTTP`服务器会自动处理此请求处理程序不需要此字段。
* 对于客户端请求，设置此字段可防止重复使用对同一主机的请求之间的`TCP`连接，
* 如同设置了`Transport.DisableKeepAlives`。
#### 主机 `Host string`
* 对于服务器请求Host指定其上的主机寻找网址。根据RFC 2616，header host键的值或URL本身中给出的主机名。
* 它可能是“host：port”的形式。对于国际域名名称，主机可以是Punycode或Unicode形式。
* 使用 golang.org/x/net/idna将其转换为任何格式如果需要,对于客户端请求主机可选择覆盖主机 要发送的标头。 
* 如果为空，则Request.Write方法使用URL.Host的值。 Host 可能包含国际域名。
#### 请数据 `Form url.Values`
* 表单包含已解析的表单数据，包括URL字段的查询参数和POST或PUT表单数据。
* 此字段仅在调用ParseForm后可用。HTTP客户端忽略Form并改为使用Body。
#### 请求体 `PostForm url.Values`
* PostForm包含来自POST，PATCH的已解析表单数据，或PUT体参数。	
* 此字段仅在调用ParseForm后可用。HTTP客户端忽略PostForm并使用Body。
#### 请求体 `MultipartForm *multipart.Form`
* MultipartForm是解析更广扩展部分表单，包括文件上载。此字段仅在调用ParseMultipartForm后可用。
* HTTP客户端忽略MultipartForm并改为使用Body。
#### `Trailer Header`
* Trailer指定在请求体后发送的其headers。
* 对于服务器请求，Trailer map 最初只包含 Trailer密钥，nil值。
* （客户声明哪些Trailer稍后会发送。）当处理程序从Body读取时，它必须 不是参考Trailer。
* 读取Body后返回EOF，Trailer可以由客户端再次读取并包含非零值（如果已发送）。
* 对于客户端请求，必须将Trailer初始化为map包含的以后发送的Trailer键。 
* 值可以是零或最终值。 ContentLength必须为0或-1，以发送分块请求。发送HTTP请求后，可以更新map值 读取请求正文。 
* 一旦body返回EOF，调用者则者必须 不要改变Trailer。很少有HTTP客户端，服务器或代理支持HTTP Trailer。
#### `RemoteAddr string`
* RemoteAddr允许HTTP服务器和其他软件进行记录发送请求的网络地址，通常用于记录。
* ReadRequest并未填写此字段没有定义的格式。 
* 此包中的HTTP服务器在调用之前将RemoteAddr设置为“IP：port”地址处理程序。HTTP客户端忽略此字段。
#### 请求体 `RequestURI string`
* RequestURI是未经修改的Request-URI,客户端发送的Request-Line（RFC 2616，第5.1节）
* 到服务器 通常应该使用URL字段。在HTTP客户端请求中设置此字段是错误的。
#### `TLS *tls.ConnectionState`
* TLS允许HTTP服务器和其他软件进行记录有关请求的TLS连接的信息。 ReadRequest未填写此字段。
* 此程序包中的HTTP服务器为其设置字段
* 在调用处理程序之前启用TLS连接;否则它会使字段为零。HTTP客户端忽略此字段。
#### 请求体 `Cancel <-chan struct{}`
* Cancel是一个可选通道，其闭包表示客户端请求应被视为已取消。 
* 不是所有的实现 RoundTripper 支持取消。对于服务器请求，此字段不适用。
* 不推荐使用：使用Context和WithContext方法来替代。未定义将不会有作用，是否遵守取消。
#### 请求体 `Response *Response`
* 此字段仅在客户端重定向期间填充。
#### 请求体 `ctx context.Context`
* ctx是客户端或服务器上下文。它应该只通过使用WithContext复制整个请求来修改。
* 为了防止人们使用Context应为错误是不可预知的,改变同一请求的调用者持有的上下文。