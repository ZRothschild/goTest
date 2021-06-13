package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"time"
)

type S struct {
	Name string
	Msg  chan string
}

func (s *S) PrintName() {
	go func() {
		fmt.Println("####")
		fmt.Println(<-s.Msg)
		fmt.Println("999999999")
	}()
	// fmt.Println("0000")
}

func (s *S) SetName() {
	// go func() {
	fmt.Println("===")
	fmt.Println("--------")
	s.Msg <- "弄好了"
	// close(s.Msg)
	fmt.Println("xxxxxxxx")
	// }()
	fmt.Println(1111)
	return
}

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	fmt.Println("===")
	testOne()
	fmt.Println("=========出来了")
	s := <-c
	fmt.Println("Got signal:", s)
	return
	// s := S{Msg: make(chan string)}
	// s.SetName()
	// time.Sleep(4*time.Second)
	// s.PrintName()
	// time.Sleep(4*time.Second)
	// arr := []int{1, 2, 3, 4}
	// for _, v := range arr {
	// 	arr = append(arr, v)
	// }
	// fmt.Println(arr)
	// var mapSet = sync.Map{}
	// mapSet.Store("now","test")
	// mapSet.Delete("now")
	// mapSet.LoadOrStore("zhao","zhao")

	ch := make(chan int)
	go func() {
		//fmt.Println("Hello inline")
		ch <- 1
		//fmt.Println("Hello ")
		j := 2
		j++
	}()
	go printHello(ch)

	time.Sleep(2)
	//fmt.Println("Hello from main")
	<-ch
	//fmt.Println("Recieved i", i)
	<-ch
	//fmt.Println("Recieved j", j)


	fmt.Println("Commencing	countdown.")
	tick := time.Tick(1 * time.Second)
	fmt.Println(tick)
	for count := 10; count > 0; count-- {
		fmt.Println(count)
		//i := <-tick
		//fmt.Println(i)
	}
	fmt.Println("ggggggggg")
	//
	t, err := time.Parse("2006-01-02 15:04:05", "2018-09-30 10:48:20")

	fmt.Println(err)
	fmt.Println("ggggggggg")

	aa := t.AddDate(0, 0, 1).Format("2006-01-02")
	fmt.Println(aa)
	//t,_ := time.ParseInLocation("2006-01-02T","2018-09-30T10:48:20Z",time.Local)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	a := time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05")
	fmt.Println(a)

}

func testOne() {
	go test()
	return
}

func test() {
	go func() {
		for {
			select {
			default:
				time.Sleep(3 * time.Second)
				fmt.Println("shui jiao")
				time.Sleep(2 * time.Second)
				fmt.Println("shui jiao11111")
				return
			}
		}
	}()
	// time.Sleep(time.Second)
	// b <- true

	fmt.Println("aaaaaa")
	runtime.Goexit() //超时后退出该Go协程
	return
}


func printHello(ch chan int) {
	//fmt.Println("Hello from printHello111")
	ch <- 2
	//fmt.Println("Hello from printHello2222")
	i := 1
	i++
}