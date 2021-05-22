package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	file, err := os.Open("./file/text.log")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	//读取文件
	scanner := bufio.NewScanner(file)
	//循环行
	for scanner.Scan() {
		//行内容
		fmt.Println(scanner.Text())
	}
	//错误处理
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
