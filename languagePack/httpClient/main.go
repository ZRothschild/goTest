package main

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// 你需要先去了解下 golang中   & |  <<  >> 几个符号产生的运算意义
const (
	workerBits  uint8 = 10 // 10bit工作机器的id，如果你发现1024台机器不够那就调大次值
	numberBits  uint8 = 12 // 12bit 工作序号，如果你发现1毫秒并发生成4096个唯一id不够请调大次值
	workerMax   int64 = -1 ^ (-1 << workerBits)
	numberMax   int64 = -1 ^ (-1 << numberBits)
	timeShift   uint8 = workerBits + numberBits
	workerShift uint8 = numberBits
	// 如果在程序跑了一段时间修改了epoch这个值 可能会导致生成相同的ID，
	// 这个值请自行设置为你系统准备上线前的精确到毫秒级别的时间戳，因为雪花时间戳保证唯一的部分最多管69年（2的41次方），
	// 所以此值设置为你当前时间戳能够保证你的系统是从当前时间开始往后推69年
	startTime int64 = 1525705533000
)

type Worker struct {
	mu        sync.Mutex
	timestamp int64
	workerId  int64
	number    int64
}

func NewWorker(workerId int64) (*Worker, error) {
	if workerId < 0 || workerId > workerMax {
		return nil, errors.New("Worker ID excess of quantity")
	}
	// 生成一个新节点
	return &Worker{
		timestamp: 0,
		workerId:  workerId,
		number:    0,
	}, nil
}

func (w *Worker) GetId() int64 {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now().UnixNano() / 1e6
	if w.timestamp == now {
		w.number++
		if w.number > numberMax {
			for now <= w.timestamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		w.number = 0
		w.timestamp = now
	}
	// 以下表达式才是主菜
	//  (now-startTime)<<timeShift   产生了 41 + （10 + 12）的效应但却并不保证唯一
	//  | (w.workerId << workerShift)  保证了与其他机器不重复
	//  | (w.number))  保证了自己这台机不会重复
	ID := int64((now-startTime)<<timeShift | (w.workerId << workerShift) | (w.number))
	return ID
}

func main() {
	fmt.Printf("%s  %s", strconv.FormatInt(285439180060758016, 2), strconv.FormatInt(285439203959902208, 2))

	// 生成节点实例，当你分布式的部署你的服务的时候，这个NewWorker的参数记录不同的node配置的值应该不一样
	// node, err := NewWorker(1)
	// if err != nil {
	// 	panic(err)
	// }
	// for {
	// 	fmt.Println(node.GetId())
	// 	return
	// }
}

// package main
//
// import (
// 	// "crypto/tls"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"fmt"
// 	"net/url"
// )
//
//
//
// func main() {
// 	tr := &http.Transport{
// 		// TLSClientConfig: &tls.Config{RootCAs: pool},
// 		Proxy: func(request *http.Request) (url *url.URL, e error) {
// 			return url.Parse("http://0.0.0.0:8080")
// 		},
// 		DisableCompression: true,
// 	}
//
// 	client := &http.Client {
// 		CheckRedirect: redirectPolicyFunc,
// 		Transport: tr,
// 	}
// 	resp, err := client.Get("http://127.0.0.1:8080/test")
// 	if err != nil {
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("body %s \n", err)
// 	}
// 	fmt.Println(string(body))
//
// 	req, err := http.NewRequest("GET", "http://127.0.0.1:8080/list", nil)
// 	resp, err = client.Do(req)
//
// 	defer resp.Body.Close()
// 	body, err = ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		log.Printf("body %s \n", err)
// 	}
// 	fmt.Println(string(body))
//
// }
//
// func redirectPolicyFunc  (req *http.Request,via []*http.Request) error {
// 	fmt.Println("1111111111")
// 	fmt.Println(req.URL)
// 	for k,v :=range via{
// 		fmt.Printf("%v  %v\n",k,v)
// 	}
// 	return nil
// }
