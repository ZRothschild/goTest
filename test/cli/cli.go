package main

// docker run -p 3307:3306  --restart=always  --privileged=true --name mysql -v D:\dockerData\mysql:/var/lib/mysql -v D:\dockerData\mysql\my.cnf:/etc/mysql/my.cnf -e MYSQL_ROOT_PASSWORD="123456" -d mysql:5.7

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"grpcTest/mes"
	"grpcTest/name"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":5001", grpc.WithBlock(), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("请求服务错误 %s\n", err)
		return
	}
	defer conn.Close()

	waiterCli := mes.NewWaiterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := waiterCli.HelloTest(ctx, &mes.Req{Name: "李子明", NN: &name.Res{
		Name: "留一手",
	}})

	fmt.Printf("数据 %#v\n", res)

	if err != nil {
		//fmt.Println(err)
		log.Fatalf("请求方法错误 %s \n", err.Error())
		return
	}

	aa, err := json.Marshal(res)

	fmt.Printf("服务返回 %s", string(aa))
}
