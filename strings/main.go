package main

import (
	"fmt"
	"time"
)

func main()  {
	const n = 60
	starttime := time.Now()
	fibN := fib(n)
	endtime := time.Now()
	cost_time := endtime.Sub(starttime)
	fmt.Println(cost_time)
	fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}