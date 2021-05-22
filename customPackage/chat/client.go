package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

var user map[string]string

func main() {
	//tcp 连接8080端口
	conn, err := net.Dial("tcp", ":8080")
	defer conn.Close()
	if err != nil {
		fmt.Printf("net.Dial %v\n", err)
		os.Exit(1)
	}
	//发送信息向服务端
	go sendMessage(conn)
	buf := make([]byte, 1024)
	for {
		len, _ := conn.Read(buf)
		if len == 0 {
			fmt.Printf("net.Dial conn.Read %v\n", err)
			break
		}
		fmt.Printf("server send %s\n", string(buf[:len]))
	}
	fmt.Printf("aaaa  %s\n", "成功")
}

//发送消息
func sendMessage(conn net.Conn) {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data, _, _ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}
		_, err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Printf("sendMessage Write %s\n", err.Error())
			break
		}
	}
}

func readMessage(conn net.Conn) {
	buf := make([]byte, 1024)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("net.Dial conn.Read %v\n", err)
			break
		}
		fmt.Printf("server send %s\n", string(buf))
	}
}
