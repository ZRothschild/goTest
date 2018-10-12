package main

import (
	"fmt"
)

func printHello(ch chan int) {
	fmt.Println("Hello from printHello111")
	ch <-2
	//fmt.Println("Recieved i", i)
	fmt.Println("Hello from printHello2222")
}
func main() {

	ch := make(chan int)
	go func() {
		fmt.Println("Hello inline")
		ch <- 1
		fmt.Println("Hello ")
	}()
	go printHello(ch)
	fmt.Println("Hello from main")
	i := <-ch

	fmt.Println("Recieved i", i)
	 <-ch
	fmt.Println("Recieved j", 0)
}