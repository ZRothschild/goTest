package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	//multipart/form-data
	//application/x-www-form-urlencoded
	//application/json

	buf := new(bytes.Buffer)
	wr := multipart.NewWriter(buf)
	defer wr.Close()
	wr.WriteField("lang", "go")
	wr.WriteField("web", "iris")

	filePath := "/home/zrothschild/testWeb/main.go"
	cf, _ := wr.CreateFormFile("fileName", filePath)

	of, _ := os.Open(filePath)
	defer of.Close()
	io.Copy(cf, of)

	ct := wr.FormDataContentType()
	resp, err := http.Post("http://127.0.0.1:8000/list", ct, buf)

	if err != nil {
		log.Printf("resp %s \n", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("body %s \n", err)
	}
	fmt.Println(string(body))
}
