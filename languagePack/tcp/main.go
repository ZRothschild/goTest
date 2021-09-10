package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	//conn, err := net.Dial("tcp", "baidu.com:80")
	//if err != nil {
	//	// handle error
	//}
	//fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	//status, err := bufio.NewReader(conn).ReadString('\n')
	//fmt.Println(status)

	la, err := net.ResolveUDPAddr("udp4", "127.0.0.1:0")
	if err != nil {
		fmt.Println(err)
	}
	c, err := net.ListenUDP("udp4", la)
	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	c.LocalAddr()
	c.RemoteAddr()
	c.SetDeadline(time.Now().Add(10))
	c.SetReadDeadline(time.Now().Add(10))
	c.SetWriteDeadline(time.Now().Add(10))
	c.SetReadBuffer(2048)
	c.SetWriteBuffer(2048)



	wb := []byte("UDPCONN TEST")
	rb := make([]byte, 128)
	if _, err := c.WriteToUDP(wb, c.LocalAddr().(*net.UDPAddr)); err != nil {
		fmt.Println(err)
	}
	if _, _, err := c.ReadFromUDP(rb); err != nil {
		fmt.Println(err)
	}
	if _, _, err := c.WriteMsgUDP(wb, nil, c.LocalAddr().(*net.UDPAddr)); err != nil {
		fmt.Println(err)
	}
	if _, _, _, _, err := c.ReadMsgUDP(rb, nil); err != nil {
		fmt.Println(err)
	}

	if f, err := c.File(); err != nil {
		fmt.Println(err)
	} else {
		f.Close()
	}

	defer func() {
		if p := recover(); p != nil {
			fmt.Println(err)
		}
	}()

	c.WriteToUDP(wb, nil)
	c.WriteMsgUDP(wb, nil, nil)


}