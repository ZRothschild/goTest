package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type Arith int

type Arguments struct {
	A, B int
}

func (a *Arith) AddAction(arg *Arguments, resulte *int) error {
	*resulte = arg.A + arg.B
	return nil
}

func main() {
	arith := new(Arith)

	//typ := reflect.TypeOf(arith)
	//rcvr := reflect.ValueOf(arith)
	//sname := reflect.Indirect(rcvr).Type().Name()
	//fmt.Printf("%v\n%#v\n%#v\n%v\n",typ,rcvr,sname,reflect.Indirect(rcvr))

	rpc.RegisterName("Test", arith)
	rpc.HandleHTTP()
	err := http.ListenAndServe(":5555", nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	//tcpAdr,_ := net.ResolveTCPAddr("tcp",":5555")
	//listen ,_ := net.ListenTCP("tcp",tcpAdr)
	//for  {
	//	conn,err := listen.Accept()
	//	if err != nil{
	//		fmt.Printf("error %v",err)
	//		continue
	//	}
	//	rpc.ServeConn(conn)
	//}
}
