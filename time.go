package main

import (
	"time"
	"fmt"
)

func main()  {

	/* 单位换算
	 * 1秒(s)    =  1000毫秒(ms)
	 * 1毫秒(ms) =  1000微秒(μs)
	 * 1微秒(μs) =  1000纳秒(ns)
	 * 1纳秒(ns) =  1000皮秒(ps)
	 * 单位间隔为 1000
	 */

	//获取当前时间戳 精确到秒 也就是十位数字 例如：1539138946
	timestamp := time.Now().Unix()
	fmt.Printf("时间戳秒： %d\n",timestamp)

	//时间转成时间戳
	times := time.Unix(timestamp, 0).Format("2006-02-01 15:04:05")
	fmt.Printf("times：%s\n",times)

	//time.Unix 第一个参数是秒时间戳，第二个是纳秒 范围在 [0, 999999999]
	timen := time.Unix(timestamp, 500).Format("2006-02-01 15:04:05.000")
	fmt.Printf("timen：%s\n",timen)

	//获取当前时间戳 精确到纳秒 也就是19位数字 例如：1539140887816267500
	stamp := time.Now().UnixNano()
	fmt.Printf("时间戳纳秒：%d\n",stamp)

	//获取当前时间戳 精确到毫秒 也就是13位数字  例如：1539144568939
	fmt.Printf("时间戳毫秒：%d\n",time.Now().UnixNano() / 1e6)

	//时间戳转时间格式,Format接收，一个时间格式  例如：Wed Oct 10 12:09:28 2018
	ANSIC := time.Now().Format(time.ANSIC)
	fmt.Printf("ANSIC格式：%s\n",ANSIC)

	//Format函数也可以自定义格式,写一个常量
	//go比较特殊不使用 yyyyMMdd hh:mm:ss 而是 2006 02 01 15:04:05 自行组合
	formatOne := time.Now().Format("01 02 15:04:05 2006")
	fmt.Printf("1 2 15:04:05 2006 格式：%s\n",formatOne)

	//如果想要 yyyy/MM/dd 应该如下 例如：2018/10/11
	formatTwo := time.Now().Format("2006/02/01")
	fmt.Printf(" yyyy/MM/dd 格式：%s\n",formatTwo)

	//将字符时间，转成Time类型，二者类型应该保持一致
	checkTime,_ := time.Parse("2006-01-02T15:04:05Z","2018-10-10T12:25:30Z")
	fmt.Printf("checkTime %s\n",checkTime.Format("2006-01-02 15:04:05"))
}