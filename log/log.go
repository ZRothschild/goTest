package main

import (
	"log"
	"os"
)

func main() {
	//O_APPEND 添加写  O_CREATE 不存在则生成   O_WRONLY 只写模式
	f, err := os.OpenFile("./log/text.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	//时间展示格式  LstdFlags
	logger := log.New(f, "prefix ", log.LstdFlags)
	logger.Println("text to append")
	logger.Println("more text to append")
}
