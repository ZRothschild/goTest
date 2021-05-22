package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

var (
	//nsqd的地址，使用了tcp监听的端口
	tcpNsqdAddrr = "127.0.01:4150"
)

//声明一个结构体，实现HandleMessage接口方法（根据文档的要求）
type NsqHandler struct {
	//消息数
	msqCount int64
	//标识ID
	nsqHandlerID string
}

//实现HandleMessage方法
//message是接收到的消息
func (s *NsqHandler) HandleMessage(message *nsq.Message) error {
	//没收到一条消息+1
	s.msqCount++
	//打印输出信息和ID
	fmt.Println(s.msqCount, s.nsqHandlerID, time.Now())
	//打印消息的一些基本信息
	fmt.Printf(" test msg.Timestamp=%v, msg.nsqaddress=%s,msg.body=%s \n", time.Unix(0, message.Timestamp).Format("2006-01-02 03:04:05"), message.NSQDAddress, string(message.Body))
	return nil
}

func main() {
	//初始化配置
	//config := nsq.NewConfig()
	//for i := 0; i < 50; i++ {
	//	//创建100个生产者
	//	tPro, err := nsq.NewProducer(tcpNsqdAddrr, config)
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//	//主题
	//	topic := "Insert"
	//	//主题内容
	//	tCommand := "new data!"
	//	//发布消息
	//	//err = tPro.Publish(topic, []byte(tCommand))
	//	err = tPro.DeferredPublish(topic, 1 * time.Minute,[]byte(tCommand+"==="+strconv.Itoa(i)+"==="+time.Now().Format("2006/1/2 15:04:05")))
	//	if err != nil {
	//		fmt.Println(err)
	//	}
	//}

	//初始化配置
	config := nsq.NewConfig()
	//创造消费者，参数一时订阅的主题，参数二是使用的通道
	com, err := nsq.NewConsumer("Insert", "channel1", config)
	if err != nil {
		fmt.Println(err)
	}
	//添加处理回调
	com.AddHandler(&NsqHandler{nsqHandlerID: "One"})
	//连接对应的nsqd
	err = com.ConnectToNSQD(tcpNsqdAddrr)
	if err != nil {
		fmt.Println(err)
	}

	time.Sleep(time.Millisecond)
	//只是为了不结束此进程，这里没有意义
	var wg = &sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
}
