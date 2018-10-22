package main

import (
	"net"
	"os"
	"fmt"
	"bufio"
	"strings"
)

func main()  {
	conn ,err := net.Dial("tcp",":8080")
	defer conn.Close()
	if err != nil {
		fmt.Printf("net.Dial %v\n",err)
		os.Exit(1)
	}
	go sendMessage(conn)
	buf := make([]byte,1024)
	for {
		_,err := conn.Read(buf)
		if err != nil {
			fmt.Printf("net.Dial conn.Read %v\n",err)
			break
		}
		fmt.Printf("server send %s\n",string(buf))
	}
	fmt.Printf( "aaaa  %s\n","成功")
}

func sendMessage(conn net.Conn)  {
	var input string
	for {
		reader := bufio.NewReader(os.Stdin)
		data ,_,_ := reader.ReadLine()
		input = string(data)

		if strings.ToUpper(input) == "EXIT" {
			conn.Close()
			break
		}
		_,err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Printf("sendMessage Write %s\n",err.Error())
			break
		}
	}
}

func readMessage(conn net.Conn)  {
	buf := make([]byte,1024)
	for  {
		_,err := conn.Read(buf)
		if err != nil {
			fmt.Printf("net.Dial conn.Read %v\n",err)
			break
		}
		fmt.Printf("server send %s\n",string(buf))
	}
}