package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"grpcTest/mes"
	"grpcTest/waiter"
	"net/http"
	"strings"

	//"net"
)

func main() {
	//l, err := net.Listen("tcp", "localhost:5001")
	//if err != nil {
	//	log.Fatalf("监听错误 %s\n", err)
	//	return
	//}
	mux := runtime.NewServeMux()

	//s := grpc.NewServer()
	//mes.RegisterWaiterServer(s, &waiter.Waiter{})
	if err := mes.RegisterWaiterHandlerServer(context.Background(), mux, &waiter.Waiter{}); err != nil {
		return
	}

	if err := http.ListenAndServe(
		":5001",
		mux,
		//"./x509/server_cert.pem",
		//"./x509/server_key.pem",
		//grpcHandlerFunc(s, mux),
		//http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//	log.Printf("收到请求%v", r)
		//	s.ServeHTTP(w, r)
		//}),
	); err != nil {
		fmt.Println(err)
		return
	}
	//if err = s.Serve(l); err != nil {
	//	log.Fatalf("服务启动错误 %s\n", err)
	//	return
	//}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	fmt.Println("aaaaaaaa")
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
