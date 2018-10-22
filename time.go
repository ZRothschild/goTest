package main

import (
	"time"
	"fmt"
)

func main()  {

	/**
	 *单位换算
	 * 1秒(s)    =  1000毫秒(ms)
	 * 1毫秒(ms) =  1000微秒(μs)
	 * 1微秒(μs) =  1000纳秒(ns)
	 * 1纳秒(ns) =  1000皮秒(ps)
	 * 单位间隔为 1000
	 *
	 */

	//获取当前时间戳 精确到秒 也就是十位数字 例如：1539138946
	timestamp := time.Now().Unix()
	fmt.Printf("时间戳秒： %d\n",timestamp)

	//时间转成时间戳
	times := time.Unix(timestamp, 0).Format("2006-02-01 15:04:05")
	fmt.Printf("times：%s\n",times)

	//time.Unix 第一个参数是秒时间戳，第二个是纳秒 范围在 [0, 999999999]
	timen := time.Unix(timestamp, 0).Format("2006-02-01 15:04:05.000")
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
	checkOne,_ := time.Parse("2006-01-02T15:04:05Z","2018-10-10T12:25:30Z")
	fmt.Printf("checkOne %s\n",checkOne.Format("2006-01-02 15:04:05"))

	//如果你的字符串时间是 2018-10-10 12:25:30 这应该如下转换
	checkTwo,_ := time.Parse("2006-01-02 15:04:05","2018-10-10 12:25:30")
	fmt.Printf("checkTwo %s\n",checkTwo.Format("2006/01/02 15:04:05"))
	//checkTwo 2018/10/10 12:25:30

	//因为Parse使用默认UTC国际默认时区，time.Unix使用的本地时区，中国为 +8 东八时区，所以会自动加上八小时
	checkThree := time.Unix(checkTwo.Unix(),0).Format("2006/01/02 15:04:05")
	fmt.Printf("checkThree %s\n",checkThree)   //checkThree 2018/10/10 20:25:30

	/**
	   时间显示会有与时区相关函数
	   func (t Time) UTC() Time {} //以UTC 时区的时间表示
	   func (t Time) Local() Time {} // 以本地时区表示
	   func (t Time) In(loc *Location) Time {} // 按指时区式显示时间
	 */
	//解决这个问题当然可以将去八个小时的秒数,也可以用另外一种方法
	parLoc,_ := time.ParseInLocation("2006-01-02 15:04:05","2018-10-10 12:25:30",time.Local)
	fmt.Printf("Local %s\n",parLoc.Local()) //2018-10-10 12:25:30 +0800 CST
	fmt.Printf("UTC %s\n",parLoc.UTC()) //UTC 2018-10-10 04:25:30 +0000 UTC
	fmt.Printf("In %s\n",parLoc.In(time.Local))//In 2018-10-10 12:25:30 +0800 CST

	parTime := time.Unix(parLoc.Unix(),0).Format("2006/01/02 15:04:05")
	fmt.Printf("parTime %s\n",parTime)   //parTime 2018/10/10 12:25:30

	/**
		func (t Time) Date() (year int, month Month, day int) {} // 返回年月日
		func (t Time) Year() int {} // 返回年
		func (t Time) Month() Month {} // 月
		func (t Time) Day() int {} // 日
		func (t Time) Weekday() Weekday {} // 星期
		func (t Time) ISOWeek() (year, week int) {} // 返回年，星期范围编号
		func (t Time) Clock() (hour, min, sec int) {} // 返回时间的时分秒
		func (t Time) Hour() int {} // 返回小时
		func (t Time) Minute() int {} // 分钟
		func (t Time) Second() int {} // 秒
		func (t Time) Nanosecond() int {} // 纳秒
		func (t Time) YearDay() int {} // 一年中的天第几天
		func (t Time) Location() *Location {} //时区
		func (t Time) Zone() (name string, offset int) {} // 时间所在时区与UTC时区的时差
	 */

	 //返回年月日
	y,m,d := time.Now().Date()
	fmt.Printf("%d 年 %d 月 %d 日\n",y,m,d)  //2018 年 10 月 17 日
	//返回时分秒
	h,mi,s := time.Now().Clock()
	fmt.Printf("%d 时 %d 分 %d 秒\n",h,mi,s) //16 时 21 分 56 秒

	//一年第多少天
	fmt.Printf("一年的第%d天\n",time.Now().YearDay()) //一年的第290天

	//时区和相差的秒数
	name,off := time.Now().Zone()
	fmt.Printf("local %s stamp %d\n",name,off) //local CST stamp 28800

	//时区和相差的秒数
	y,w := time.Now().ISOWeek()
	fmt.Printf("%d 年第 %d 个星期\n",y,w) //2018 年第 43 个星期

	//func (t Time) IsZero() bool {} // 是否是零时时间
	//func (t Time) After(u Time) bool {} // 时间在u 之前
	//func (t Time) Before(u Time) bool {} // 时间在u 之后
	//func (t Time) Equal(u Time) bool {} // 时间与u 相同
	//func (t Time) Add(d Duration) Time {} // 返回t +d 的时间点
	//func (t Time) Sub(u Time) Duration {} // 返回 t-u
	//func (t Time) AddDate(years int, months int, days int) Time {} 返回增加了给出的年份

	//是否是零时时间
	isZero := time.Now().IsZero()
	fmt.Printf("是否零时时间 %t\n",isZero) //是否零时时间 false




}