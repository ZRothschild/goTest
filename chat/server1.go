package main

import (
	"net"
	"fmt"
	"os"
)

func main()  {
	listener , err := net.Listen("tcp",":8080")
	defer listener.Close()
	if err != nil {
		fmt.Printf("listener err %v\n",err)
		os.Exit(1)
	}
	for  {
		conn,err := listener.Accept()
		go process(conn,err)
	}
}

func process(conn net.Conn,err error)  {
	if err != nil {
		fmt.Printf("process err %v\n",err)
		os.Exit(1)
	}
	buf := make([]byte,1024)
	defer conn.Close()
	for  {
		//把信息读取到buf 里面
		_ , err = conn.Read(buf)
		if err != nil {
			break
		}
		fmt.Printf("message %s\n",string(buf))
	}
}