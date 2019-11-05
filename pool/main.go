package main

import (
	"fmt"
	"sync"
	"unsafe"
)

func main() {
	p := &sync.Pool{
		New: func() interface{} {
			return 0
		},
	}
	a := p.Get().(int)
	p.Put(1)
	b := p.Get().(int)
	fmt.Println(a, b)

	num := 5
	numPointer := &num

	flnum := (*float32)(unsafe.Pointer(numPointer))
	fmt.Println(flnum)
}
