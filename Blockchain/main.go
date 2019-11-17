package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

func main() {
	sha1New := md5.New()
	src := "123456"
	srcB, _ := hex.DecodeString(src)
	fmt.Printf("123456 r => %x\n", srcB)
	b := []byte("123456")
	sha1New.Write(b)
	nSum := sha1New.Sum(nil)
	fmt.Printf("123456 r => %x\n", nSum)
	sha1New.Reset()
	sha1New.Write(srcB)
	nSum = sha1New.Sum(nil)
	fmt.Printf("123 r => %x\n", nSum)
}
