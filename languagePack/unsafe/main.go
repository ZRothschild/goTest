package main

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"unsafe"
)

func main() {
	bt := []byte{213,212,199,197,205,250}



	fmt.Println(bt)
	var t int = 3
	tP := *(*int)(unsafe.Pointer(&t))
	fmt.Printf("%d\n", tP)
}


func UTF82GB2312(s []byte)([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.HZGB2312.NewEncoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}