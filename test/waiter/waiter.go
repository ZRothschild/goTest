package waiter

import (
	"context"
	"fmt"
	"grpcTest/mes"
	"time"
)

func NewWaiterServer() mes.WaiterServer {
	return new(Waiter)
}

type Waiter struct {
	mes.UnimplementedWaiterServer
}

func (w Waiter) HelloTest(ctx context.Context, req *mes.Req) (res *mes.Res, err error) {
	af := time.After(5 * time.Second)
	//time.Sleep(2 * time.Second)
	select {
	case <-ctx.Done():
		res = &mes.Res{
			Name: "我是超时了",
		}
		return res, nil
	case <-af:
		res = &mes.Res{
			Name: "我是五秒后超时了",
		}
		return res, nil
	default:
		fmt.Printf("ni hao %v\n", req)
		res = &mes.Res{
			Name: req.Name,
			//Img: []*mes.Req{},
		}
		return res, nil
	}
}

func (w Waiter) Echo(ctx context.Context, req *mes.Req) (res *mes.Res, err error) {
	af := time.After(5 * time.Second)
	time.Sleep(2 * time.Second)
	select {
	case <-ctx.Done():
		res = &mes.Res{
			Name: "我是超时了",
		}
		return res, nil
	case <-af:
		res = &mes.Res{
			Name: "我是五秒后超时了",
		}
		return res, nil
	default:
		fmt.Printf("ni hao %v\n", req)
		res = &mes.Res{
			Name: req.Name,
		}
		return res, nil
	}
}
