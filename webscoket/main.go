package main

import (
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"reflect"
)

func test(ws *websocket.Conn) {
	data := map[string]interface{}{}
	for {
		err := websocket.JSON.Receive(ws, &data)
		if err != nil {
			fmt.Printf("err %v", err)
			ws.Close()
			break
		}
		fmt.Printf("Server receive : %v\n", data)
		fmt.Printf("Type received : %s - %s - %s\n", reflect.TypeOf(data["texts"]), reflect.TypeOf(data["flage"]), reflect.TypeOf(data["bool"]))
		data["texts"] = data["texts"]
		data["flage"] = 234
		data["bool"] = false
		err2 := websocket.JSON.Send(ws, data)
		if err2 != nil {
			fmt.Println(err2)
			break
		}
	}
}

func main() {
	http.Handle("/test/", websocket.Handler(test))
	err := http.ListenAndServe(":5555", nil)
	fmt.Printf("http %s\n", err)
}
