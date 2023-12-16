package main

import "fmt"

func main() {
	defer func() {
		fmt.Println("defer 0")
		if err := recover(); err != nil {
			fmt.Println("defer recover", err)
		}
	}()

	defer func() {
		fmt.Println("defer 1")
	}()

	test()

	//panic("panic wow!")

	fmt.Println("hello world")
}

func test() {
	//defer func() {
	//	fmt.Println("test defer 0")
	//	if err := recover(); err != nil {
	//		fmt.Println("test defer recover", err)
	//	}
	//}()
	defer func() {
		fmt.Println("test defer 1")
	}()

	panic("test panic wow!")

	fmt.Println("test hello world\n")
}
