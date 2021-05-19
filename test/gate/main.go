package main

import (
	"context"
	"fmt"
	"github.com/golang/glog"

	//"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"grpcTest/mes"
	"grpcTest/waiter"
	"net/http"
	"strings"
)

func main() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err := mes.RegisterWaiterHandlerFromEndpoint(ctx, mux, "localhost:5001", opts)
	if err != nil {
		fmt.Println(err)
		return
	}
	s := grpc.NewServer()
	mes.RegisterWaiterServer(s, &waiter.Waiter{})
	//s := &http.Server{
	//	Addr:    ":5001",
	//	Handler: mux,
	//}
	//go func() {
	//	<-cxt.Done()
	//	glog.Infof("Shutting down the http gateway server")
	//	if err := s.Shutdown(context.Background()); err != nil {
	//		glog.Errorf("Failed to shutdown http gateway server: %v", err)
	//	}
	//}()
	//defer glog.Flush()

	//if err = http.ListenAndServeTLS(
	//	":5001",
	//	"./x509/server_cert.pem",
	//	"./x509/server_key.pem",
	//	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		log.Printf("收到请求%v", r)
	//		s.ServeHTTP(w, r)
	//	}),
	//); err != nil {
	//	fmt.Println(err)
	//	return
	//}
	// 如果不加密使用 http2
	if err := http.ListenAndServe(":5001", mux); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
	}
}

func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
