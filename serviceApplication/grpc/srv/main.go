package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"serviceApplication/grpc/balancer"
	pb "serviceApplication/grpc/helloworld"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const svcName = "project/test"

var addr = "127.0.0.1:50051"
var etcdAddr = "127.0.0.1:2379"

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}


func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}
	defer lis.Close()
	s := grpc.NewServer()
	defer s.GracefulStop()
	pb.RegisterGreeterServer(s, &server{})
	go balancer.Register(etcdAddr, svcName, addr, 100)
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		s := <-ch
		balancer.UnRegister(svcName, addr)
		if i, ok := s.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}