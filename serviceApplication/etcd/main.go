package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	//客户端配置
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}
	//建立连接
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("connect success")
	defer client.Close()

	//控制超时
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//1. 增-存值
	putRepose, err := client.Put(ctx, "demo1_key", "demo1_value")
	//操作完毕，cancel掉
	cancel()
	if err != nil {
		fmt.Println("put failed, err: ", err, putRepose)
		return
	}

	//2. 查-获取值， 也设置超时
	ctx, cancel = context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "demo1_key")
	// Get查询还可以增加WithPrefix选项，获取某个目录下的所有子元素
	//eg: resp, err := client.Get(ctx, "demo“, clientv3.WithPrefix())
	cancel()
	if err != nil {
		fmt.Println("get failed err:", err)
		return
	}
	for _, item := range resp.Kvs { //Kvs 返回key的列表
		fmt.Printf("%s : %s \n", item.Key, item.Value)
	}
	//3. 改-修改值
	ctx, _ = context.WithTimeout(context.Background(), time.Second)
	putRepose, err = client.Put(ctx, "demo1_key", "update_value", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println("get failed err: ", err)
	}
	fmt.Println(string(putRepose.PrevKv.Value))

	//4. 删-删除值
	//ctx, _ = context.WithTimeout(context.Background(), time.Second)
	//deleteResponse, err := client.Delete(ctx, "demo1_key")
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(deleteResponse.PrevKvs)
}
