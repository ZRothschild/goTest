package test

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	if err != nil {
		fmt.Printf("listener err %v\n", err)
		os.Exit(1)
	}
	for {
		conn, err := listener.Accept()
		defer conn.Close()
		if err != nil {
			fmt.Printf("process err %v\n", err)
			os.Exit(1)
		}
		go process(conn)
	}
}

func process(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		//把信息读取到buf 里面
		leng, err := conn.Read(buf)
		if err != nil {
			break
		}
		ip := conn.RemoteAddr()
		fmt.Printf("remote ip %s\n", ip)
		fmt.Printf("message %s\n", string(buf[:leng]))
	}
}
