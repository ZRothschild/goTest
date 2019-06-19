package main

import (
	"fmt"
)

func main() {

	a := []byte("你好中国")

	b := []string{"aa"}

	fmt.Printf("%v\n", a)
	fmt.Printf("%v\n", b)

	// const n = 60
	// starttime := time.Now()
	// fibN := fib(n)
	// endtime := time.Now()
	// cost_time := endtime.Sub(starttime)
	// fmt.Println(cost_time)
	// fmt.Printf("\rFibonacci(%d) = %d\n", n, fibN)
}

func fib(x int) int {
	if x < 2 {
		return x
	}
	return fib(x-1) + fib(x-2)
}
