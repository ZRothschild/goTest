package main

import (
	"net/rpc"
	"net"
	"fmt"
)

type Arith int

type Arguments struct {
	A,B int
}

func (a *Arith ) AddAction(arg *Arguments,resulte *int) error {
	*resulte = arg.A+arg.B
	return nil
}

func main()  {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	tcpAdr,_ := net.ResolveTCPAddr("tcp",":5555")
	listen ,_ := net.ListenTCP("tcp",tcpAdr)
	for  {
		conn,err := listen.Accept()
		if err != nil{
			fmt.Printf("error %v",err)
			continue
		}
		rpc.ServeConn(conn)
	}
}
