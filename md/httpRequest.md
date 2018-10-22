## Go语言模拟http发起Post,Get请求

### http包提供了HTTP客户端和服务端的实现

#### http.Get函数发出`HTTP/HTTPS`请求

>  1.`http.Get`直接拼接字符串

```go
package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
)

func main()  {
	//发出请求 
	resp, err := http.Get("https://studyiris.com/test?lang=go&web=iris")
	if err != nil {
		log.Printf("resp %s \n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n",err)
	}
	fmt.Println(string(body))
	//iris web 接口响应获取请求值
	//Context.URLParam()  Context.FormValue()
}
```

>  2.`http.Get` 通过 `net/url`包自己组装

```go
package main

import (
	"fmt"
	"net/url"
	"net/http"
	"io/ioutil"
	"log"
)

func main() {
	u, _ := url.Parse("https://studyiris.com/doc")
	q := u.Query()
	q.Set("lang", "go")
	q.Set("web", "iris")
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String());
	if err != nil {
		log.Fatal(err)
	}
	result, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", result)
    //iris web 接口响应获取请求值
    //Context.PostValue()  Context.FormValue()
}
```

#### http.Post函数发出`HTTP/HTTPS`请求

>  1.`http.Post` 提交数据格式 `application/x-www-form-urlencoded`

```go
package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"strings"
)

func main(){
	//post 请求  strings.NewReader 组装参数
	str := strings.NewReader("lang=go&web=iris")
	resp, err := http.Post("https://studyiris.com/doc","application/x-www-form-urlencoded",str)
	if err != nil {
		log.Printf("resp %s \n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n",err)
	}
	fmt.Println(string(body))
    //iris web 接口响应获取请求值
    //Context.PostValue()  Context.FormValue() c.Ctx.ReadForm(&b)
}
```

>  2.`http.Post` 提交数据格式 `application/json`

```go
package main

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	"bytes"
)

func main(){
	data := map[string]string{"lang":"go","web":"iris"}
	//把map 转换成 存放json的 []byte
	str,_ := json.Marshal(data)
	//把[]byte 转成实现了read接口的Reader结构体
	jsonData := bytes.NewReader(str)
	//发送求情
	resp, err := http.Post("https://studyiris.com/doc", "application/json;charset=utf-8",jsonData)
	if err != nil {
		log.Printf("resp %s \n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n",err)
	}
	fmt.Println(string(body))
    //iris web 接口响应获取请求值
    // type Iris struct { 
    //  	Lang  string  `form:"lang"` 
    // 	    Web  string  `form:"web"`
    // }
    //var iris Iris 
    //Context.ReadJSON(&iris)
    //Context.Writef("Received: %#+v\n", iris)
}
```

> 3.`http.Post` 提交数据格式 `multipart/form-data`

```go
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"bytes"
	"mime/multipart"
	"os"
	"io"
)

func main()  {
	//创建一个 bytes Buffer   &bytes.Buffer 
	//也可以  buf := &bytes.Buffer{}
	buf := new(bytes.Buffer)
	
	//将bytes Buffer  *Writer
	wr := multipart.NewWriter(buf)
	defer wr.Close()
	
	//往*Writer写入数据  字符串数据
	wr.WriteField("lang","go")
	wr.WriteField("web","iris")

    //文件上传
	filePath := "/home/zrothschild/testWeb/main.go"
	cf,_ := wr.CreateFormFile("fileName",filePath)
	//打开文件
	of,_ := os.Open(filePath)
	defer of.Close()
	//把文件数据复制进 cf
	io.Copy(cf,of)
	//设置请求数据类型
	ct := wr.FormDataContentType()
	resp, err := http.Post("https://studyiris.com/doc",ct,buf)

	if err != nil {
		log.Printf("resp %s \n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n",err)
	}
	fmt.Println(string(body))
//接收端 用 *http.Request.body
//返回数据格式如下
//  --4ad2fe5f62212c3023452538537eed12aedbfe925b0519190e6b63199140
//  Content-Disposition: form-data; name="lang"
//  
//    go
//  --4ad2fe5f62212c3023452538537eed12aedbfe925b0519190e6b63199140
//  Content-Disposition: form-data; name="web"
//  
//  iris
//  --4ad2fe5f62212c3023452538537eed12aedbfe925b0519190e6b63199140
//  Content-Disposition: form-data; name="fileName"; filename="/home/zrothschild/testWeb/main.go"
//  Content-Type: application/octet-stream
//  
//  package main
//  
//  import (
//  	"fmt"
//  	"net/http"
//  	"io/ioutil"
//  	"log"
//  	"bytes"
//  	"mime/multipart"
//  	"os"
//  	"io"
//  )
//  
//  func main()  {
//  	buf := new(bytes.Buffer)
//  	wr := multipart.NewWriter(buf)
//  	defer wr.Close()
//  	wr.WriteField("lang","go")
//  	wr.WriteField("web","iris")
//  
//  	filePath := "/home/zrothschild/testWeb/main.go"
//  	cf,_ := wr.CreateFormFile("fileName",filePath)
//  
//  	of,_ := os.Open(filePath)
//  	defer of.Close()
//  	io.Copy(cf,of)
//  
//  	ct := wr.FormDataContentType()
//  	resp, err := http.Post("http://127.0.0.1:8000/list",ct,buf)
//  
//  	if err != nil {
//  		log.Printf("resp %s \n",err)
//  	}
//  	defer resp.Body.Close()
//  	body, err := ioutil.ReadAll(resp.Body)
//  	if err != nil {
//  		log.Printf("body %s \n",err)
//  	}
//  	fmt.Println(string(body))
//  }
}
```

#### `http.PostForm` 函数发出`HTTP/HTTPS`请求

```go
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"net/url"
)

func main()  {
	//第二个值是 url.Values 类型
	resp,err := http.PostForm("https://studyiris.com/doc",url.Values{"lang":{"go"},"web":{"iris"}})
	if err != nil {
		log.Printf("resp %s \n",err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n",err)
	}
	fmt.Println(string(body))
    //iris web 接口响应获取请求值
    //Context.PostValue()  Context.FormValue() c.Ctx.ReadForm(&b)
}
```

#### http.NewRequest 函数发起http Post Get请求

> 1.`http.NewRequest` POST请求 提交数据格式 `application/x-www-form-urlencoded` 这个例子等同于 `http.PostForm`

```go
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"strings"
	"net/url"
)

func main()  {
	postValue := url.Values{
		"lang": {"go"},
		"web": {"iris"},
		}
	//把数据变成键值对	
	postString := postValue.Encode()
	req, err := http.NewRequest("POST","https://studyiris.com/doc", strings.NewReader(postString))
	if err != nil {
		log.Printf("err %s \n",err)
	}
	// 表单方式(必须)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//AJAX 方式请求
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("err2 %s \n",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("err1 %s \n",err)
	}
	fmt.Println(string(body))
}
```

> 2.`http.NewRequest` POST请求 提交数据格式 `application/json`

```go
package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"encoding/json"
	"bytes"
)

func main()  {
	data := map[string]string{"lang":"go","web":"iris"}
	//把map 转换成 存放json的 []byte
	str,_ := json.Marshal(data)
	//把[]byte 转成实现了read接口的Reader结构体
	jsonData := bytes.NewReader(str)
	req, err := http.NewRequest("POST","https://studyiris.com/doc", jsonData)
	if err != nil {
		log.Printf("err %s \n",err)
	}
	// 表单方式(必须)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//AJAX 方式请求
	req.Header.Add("x-requested-with", "XMLHttpRequest")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("err2 %s \n",err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("err1 %s \n",err)
	}
	fmt.Println(string(body))
}
```

> get请求，别的数据请求格式我就不一一举例了，格式都是统一的，抛砖引玉。

#### 通过`net/http`包不同请求方式（Get,Post），与不同数据请求格式总结

* `http.NewRequest`与 `http.DefaultClient.Do`是`http`请求方式`post`与`get`请求的基础`http.Get,http.Post,http.PostForm`，封装来源与他们
* 对`http.NewRequest`发起`Post,Get`请求，第三个参数为`body io.Reade`接口类型，所以只要满足请类型即可
* 应为数据请求是有格式的，所以不同请求格式需要规定真实类型实现`body io.Reade`接口即可
* `application/x-www-form-urlencoded`数据请求格式，真实类型需要字符格式实现`body io.Reade`接口，例如： `bar=baz&foo=quux`
* `application/json`数据请求格式，真实类型可以是很多种（字符串，结构体，map等）只要实现`body io.Reade`接口
* `multipart/form-data`数据请求格式,真实类型为`multipart.Writer`实现`body io.Reade`接口
* 当然对数据处理，我们还用到了很多不同包的不同方法，`json.Marshal`，`url.Values`，`ioutil.ReadAll`等等
* 我们只需要把我们的数据转换成，函数接收的类型即可。这里我们也领会到了接口的魅力
> 需要注意的是

```go
strings.NewReader(str)

func NewReader(s string) *Reader { return &Reader{s, 0, -1} }

//strings.Reader 结构体如下
//  type Reader struct {
//  	s        string
//  	i        int64 // current reading index
//  	prevRune int   // index of previous rune; or < 0
//  }

bytes.NewReader(b []byte)

func NewReader(b []byte) *Reader { return &Reader{b, 0, -1} }

//  bytes.Reader 结构体如下
//  type Reader struct {
//  	s        []byte
//  	i        int64 // current reading index
//  	prevRune int   // index of previous rune; or < 0
//  }

//  multipart.NewWriter
//  multipart.Writer 结构体如下
//  type Writer struct {
//  	w        io.Writer
//  	boundary string
//  	lastpart *part
//  }
```