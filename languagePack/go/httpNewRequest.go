package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"mime/quotedprintable"
	"net/http"
	"strings"
)

func main() {
	//创建一个 bytes Buffer   &bytes.Buffer
	//也可以  buf := &bytes.Buffer{}
	buf := new(bytes.Buffer)

	//将bytes Buffer  *Writer
	wr := multipart.NewWriter(buf)
	defer wr.Close()

	//往*Writer写入数据  字符串数据
	wr.WriteField("web", "go")
	wr.WriteField("lang", "iris")

	//设置请求数据类型
	ct := wr.FormDataContentType()

	resp, err := http.Post("http://127.0.0.1:8080/test/one", ct, buf)

	if err != nil {
		log.Printf("resp %s \n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	quotedprintable.NewReader(strings.NewReader(""))
	// multipart.NewReader(resp.Body,"---------").ReadForm(100)
	if err != nil {
		log.Printf("body %s \n", err)
	}
	fmt.Println(string(body))
}
