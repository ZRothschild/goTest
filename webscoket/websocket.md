# webscoket 小例子

### 功能介绍

后台服务，开启*webscoket*监听指定端口，前端页面输入的数据，通过`js` `websocket`发送给后台服务，
后台服务在把数据返回前端页面，前端页面接收数据，渲染到指定位置

### main.go

```go
package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"reflect"
)

func test(ws *websocket.Conn)  {
	data := map[string]interface{}{}
	for {
		err := websocket.JSON.Receive(ws,&data)
		if err != nil {
			fmt.Printf("err %v",err)
			ws.Close()
			break
		}
		fmt.Printf("Server receive : %v\n",data)
		fmt.Printf("Type received : %s - %s - %s\n",reflect.TypeOf(data["texts"]),reflect.TypeOf(data["flage"]),reflect.TypeOf(data["bool"]))
		data["texts"] = data["texts"]
		data["flage"] = 234
		data["bool"] = false
		err2 := websocket.JSON.Send(ws,data)
		if err2 != nil {
			fmt.Println(err2)
			break
		}
	}
}

func main()  {
	http.Handle("/test/",websocket.Handler(test))
	err := http.ListenAndServe(":5555",nil)
	fmt.Printf("http %s\n",err)
}
```
### index.html `es6`

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>webscoket</title>
</head>
<body>
    <h4 id="retour"></h4>
    <input type="text" id="text" />
    <button id="btn">Send</button>
    <script type="text/javascript">
        var ws = new WebSocket('ws://localhost:5555/test/');
        ws.onopen = (e) => {
            console.log('connection');
            document.getElementById('btn').addEventListener('click',(ev) => {
                e.preventDefault();
                var data = {
                    'texts': document.getElementById('text').value,
                    'flag' : 126,
                    'bool' : true,
                };
                ws.send(JSON.stringify(data));
            });
        };
        ws.onmessage = (e) => {
            var msg = JSON.parse(e.data)
            console.log(msg["texts"]);
            document.getElementById("retour").innerHTML = `${msg["texts"]} - ${msg["flage"]} - ${msg["bool"]}`;
        };
        ws.onerror = (e) => {
            console.log('error')
        };
        ws.onclose = (e) => {
            delete ws;
        };
    </script>
</body>
</html>
```

