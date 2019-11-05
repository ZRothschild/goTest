package main

import (
	"crypto/md5"
	"fmt"
)

func main() {
	sha1New := md5.New()
	b := []byte("123456")
	// a := []byte("456")
	sha1New.Write(b)
	nSum := sha1New.Sum(b)
	fmt.Printf("123456 r => %x", nSum)
}
