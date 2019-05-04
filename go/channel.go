package main

import (
	"fmt"
	//"time"
)

//prints to stdout and puts an int on channel
func printHello(ch chan int) {
	fmt.Println("Hello from printHello")
	//send a value on channel
	ch <- 2
}
func main() {
	ch := make(chan int, 2)
	go func() {
		fmt.Println("Hello inline")
		//send a value on channel
		ch <- 1
	}()
	//call a function as goroutine
	go printHello(ch)
	fmt.Println("Hello from main")
	i := <-ch
	fmt.Println("Recieved ", i)
	//time.Sleep(2*time.Second)
	close(ch)
	b := <-ch
	fmt.Println("Recievedb", b)
}
