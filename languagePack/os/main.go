package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var name string
var b = flag.Bool("b", false, "bool类型参数")
var s = flag.String("s", "", "string类型参数")

func main() {

	input := bufio.NewScanner(os.Stdin) //初始化一个扫表对象
	for input.Scan() {                  //扫描输入内容
		line := input.Text() //把输入内容转换为字符串
		fmt.Println(line)    //输出到标准输出
	}

	//   -name niHao -b true
	flag.StringVar(&name, "name", "everyone", "The greeting object.")

	//out, err := exec.Command("ipconfig", "/all").Output()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Printf("%v", string(out))

	flag.Parse()
	fmt.Println("name", name)
	fmt.Println("-b:", *b)
	fmt.Println("-s:", *s)
	fmt.Println("其他参数：", flag.Args())

	//for idx, args := range os.Args {
	//	fmt.Println("参数" + strconv.Itoa(idx) + ":", args)
	//}
}
