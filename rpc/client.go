package main

import (
	"fmt"
	"net/rpc"
)

type Arg struct {
	A, B int
}

func main() {
	client, err := rpc.DialHTTP("tcp", ":5555")
	if err != nil {
		fmt.Printf("client err %v\n", err)
	}
	args := Arg{A: 22, B: 34}
	var res int
	err = client.Call("Arith.AddAction", args, &res)
	if err != nil {
		fmt.Printf("Call err %v\n", err)
	}
	fmt.Printf("res %d\n", res)
}
