package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var t int = 3
	tP := *(*int)(unsafe.Pointer(&t))
	fmt.Printf("%d\n", tP)
}
