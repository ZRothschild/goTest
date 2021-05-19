package main

import (
	"fmt"
)

func main() {
	// sha1New := md5.New()
	// src := "123456"
	// srcB, _ := hex.DecodeString(src)
	// fmt.Printf("123456 r => %x\n", srcB)
	// b := []byte("123456")
	// sha1New.Write(b)
	// nSum := sha1New.Sum(nil)
	// fmt.Printf("123456 r => %x\n", nSum)
	// sha1New.Reset()
	// sha1New.Write(srcB)
	// nSum = sha1New.Sum(nil)
	// fmt.Printf("123 r => %x\n", nSum)
	arr := [25]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24}
	v := 856
	for _, v = range arr {
		fmt.Printf("#####%d\n ", v)
		go func() {
			// time.Sleep(1)
			fmt.Printf("=== %d\n", v)
			fmt.Printf("***** %d\n", v)
			// name(v)
		}()
	}
	fmt.Println("跑完了")
}

func name(v int) {
	fmt.Printf("===值是 %d\n", v)
}
