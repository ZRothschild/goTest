package main

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"time"

	"serviceApplication/grpc/balancer"
	pb "serviceApplication/grpc/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var etcdAddr = "127.0.0.1:2379"

func main() {
	r := balancer.NewResolver(etcdAddr)
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://zhao/project/test", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	client := pb.NewGreeterClient(conn)
	for {
		resp, err := client.SayHello(context.Background(), &pb.HelloRequest{Name: "hello"})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(resp)
		}
		<-time.After(time.Second)
	}
}