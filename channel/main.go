package main

import (
	"fmt"
	"os"
	"os/signal"
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
	s := <-c
	fmt.Println("Got signal:", s)
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
	return
}
