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
	fmt.Printf("Scheme %s\n", urlParse.Scheme)
	//输出 https

	//Opaque是否透明
	fmt.Printf("Opaque %s\n", urlParse.Opaque)
	//输出空

	//User包含了所有的认证信息，这里调用 Username和Password方法来获取独立值。
	fmt.Printf("User %s\n", urlParse.User)
	//输出 ZRothschild:123456

	//用户名
	fmt.Printf("User Username %s\n", urlParse.User.Username())
	//输出 ZRothschild

	//密码
	pass, _ := urlParse.User.Password()
	fmt.Printf("User Password %s\n", pass)
	//输出 123456

	//Host同时包括主机名和端口信息，如过端口存在的话，使用 strings.Split()
	fmt.Printf("Host %s\n", urlParse.Host)
	//输出 studyiris.com:80
	host := strings.Split(urlParse.Host, ":")
	fmt.Printf("Host %s\n", host[0])
	//输出 studyiris.com
	fmt.Printf("Port %s\n", host[1])
	//输出 80

	//路径
	fmt.Printf("Path %s\n", urlParse.Path)
	//输出 /detail

	fmt.Printf("RawPath %s\n", urlParse.RawPath)
	//输出

	fmt.Printf("ForceQuery %v\n", urlParse.ForceQuery)
	//输出 false
	//当 s = "https://ZRothschild:123456@studyiris.com:80/detail?"
	//输出 true

	fmt.Printf("RawQuery %s\n", urlParse.RawQuery)
	//输出 id=46&uid=13
	//如果想得到键值对可以
	rawMap, _ := url.ParseQuery(urlParse.RawQuery)
	fmt.Printf("rawMap %v\n", rawMap)
	//输出 map[id:[46] uid:[13]]
	//rawMap["id"][0] === 46  rawMap["uid"][0] === 13

	fmt.Printf("Fragment %s\n", urlParse.Fragment)
	//输出 jwt
}
