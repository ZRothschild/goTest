package main

import (
	"net"
	"fmt"
	"os"
	"strings"
)

var  onlineMap = make(map[string]net.Conn)
var chanQueues = make(chan string,1000)
var chanQuit = make(chan bool)

func main()  {
	listener , err := net.Listen("tcp",":8080")
	defer listener.Close()
	if err != nil {
		fmt.Printf("listener err %v\n",err)
		os.Exit(1)
	}
	go consumeMessage()
	for  {
		conn,err := listener.Accept()
		onlineMap[conn.RemoteAddr().String()] = conn
		for i := range onlineMap {
			fmt.Printf("addr %s\n",i)
		}
		defer conn.Close()
		if err != nil {
			fmt.Printf("process err %v\n",err)
			os.Exit(1)
		}
		go process(conn)
	}
}

func consumeMessage()  {
	for  {
		select {
		case msg := <- chanQueues:
			doProcessMsg(msg)
		case <-chanQuit :
			break
		}
	}
}

func doProcessMsg(msg string)  {
	content := strings.Split(msg,"#")
	if len(content) > 0{
		addr := content[0]
		msgSend := content[1]
		if  conn,ok := onlineMap[addr]; ok{
			_, err := conn.Write([]byte(msgSend))
			if err != nil {
				fmt.Println("写入错误咯 朋友")
			}
		}
	}
}

func process(conn net.Conn)  {
	buf := make([]byte,1024)
	for  {
		//把信息读取到buf 里面
		leng ,err := conn.Read(buf)
		if err != nil {
			break
		}
		message := string(buf[:leng])
		chanQueues <- message
	}
}