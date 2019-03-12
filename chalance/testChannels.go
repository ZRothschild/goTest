package main

import (
	"time"
)

func printHello(ch chan int) {
	//fmt.Println("Hello from printHello111")
	ch <- 2
	//fmt.Println("Hello from printHello2222")
	i := 1
	i++
}
func main() {
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
}
