package main

import (
	"./tj"
	"fmt"
)

func main() {
	randInt := tj.RandInt()
	var k string
	for k, _ = range randInt {
		first := randInt[k][:2]
		second := randInt[k][2:]
		b := first.Work(second)
		fmt.Println(b)
	}
}
