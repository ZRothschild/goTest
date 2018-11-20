# `go(golang)`语言使用`NoSql redis`数据库

## `send`方式

- `send`只是写在客户端缓存里面而已，并没有写入服务端数据库，需要`c.Flush()`往服务推送，`c.Receive()`是接收服务端操作返回是否成功，信号
```go
package main

import (
	"github.com/gomodule/redigo/redis"
	"fmt"
)

func main()  {
	//连接到redis数据库
	//c, err := redis.Dial("tcp", ":6379")
	//if err != nil {
	//	fmt.Printf("dial err %s\n",err)
	//}
	//defer c.Close()

	c, err := redis.DialURL("redis://localhost:6379")
	if err != nil {
		fmt.Printf("DialURL err %s\n",err)
	}
	defer c.Close()
	//选择1号数据库
	rep,err := c.Do("select",1)
	if err != nil {
		fmt.Printf("c.Do select err %s\n",err)
	}
	fmt.Printf("replay %v\n",rep)
	//set 一个key为go value值为iris的值
	err = c.Send("set","go","iris")
	if err != nil {
		fmt.Printf("set %v\n",rep)
	}
	c.Send("get","go")
	if err != nil {
		fmt.Printf("Send get %v\n",rep)
	}
	c.Flush()
	c.Receive()
	reply,err := redis.String(c.Receive())
	if err != nil {
		fmt.Printf("redis.Values %v\n",err)
	}
	fmt.Printf("get value %v\n",reply)
}
```

