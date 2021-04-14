package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Commencing	countdown.")
	tick := time.Tick(1 * time.Second)
	fmt.Println(tick)
	for count := 10; count > 0; count-- {
		fmt.Println(count)
		//i := <-tick
		//fmt.Println(i)
	}
	fmt.Println("ggggggggg")
	//
	t, err := time.Parse("2006-01-02 15:04:05", "2018-09-30 10:48:20")

	fmt.Println(err)
	fmt.Println("ggggggggg")

	aa := t.AddDate(0, 0, 1).Format("2006-01-02")
	fmt.Println(aa)
	//t,_ := time.ParseInLocation("2006-01-02T","2018-09-30T10:48:20Z",time.Local)
	fmt.Println(t.Format("2006-01-02 15:04:05"))
	a := time.Unix(t.Unix(), 0).Format("2006-01-02 15:04:05")
	fmt.Println(a)
}
