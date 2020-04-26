package main

import "unsafe"

func main()  {
	var t int
	tP := *(*int)(unsafe.Pointer(&t))
}
