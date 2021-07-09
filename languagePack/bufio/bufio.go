package main

import (
	"bufio"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"
	"time"
)

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		sum += i
	}
	return sum
}
func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (end - start + 1) * (end + start) / 2
}

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}
func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n",
				getFunctionName(f), time.Since(t))
		}(time.Now())
		return f(start, end)
	}
}

func Decorator(decoPtr, fn interface{}) (err error) {
	var decoratedFunc, targetFunc reflect.Value
	if decoPtr == nil ||
		reflect.TypeOf(decoPtr).Kind() != reflect.Ptr ||
		reflect.ValueOf(decoPtr).Elem().Kind() != reflect.Func {
		err = fmt.Errorf("Need a function porinter!")
		return
	}
	decoratedFunc = reflect.ValueOf(decoPtr).Elem()
	targetFunc = reflect.ValueOf(fn)
	if targetFunc.Kind() != reflect.Func {
		err = fmt.Errorf("Need a function!")
		return
	}
	v := reflect.MakeFunc(targetFunc.Type(),
		func(in []reflect.Value) (out []reflect.Value) {
			fmt.Println("before")
			if targetFunc.Type().IsVariadic() {
				out = targetFunc.CallSlice(in)
			} else {
				out = targetFunc.Call(in)
			}
			fmt.Println("after")
			return
		})
	decoratedFunc.Set(v)
	return
}

func bar(a, b string) string {
	fmt.Printf("%s, %s \n", a, b)
	return a + b
}

func main() {
	//mybar := bar
	//err := Decorator(&mybar, bar)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println(mybar("hello,", "world!"))

	// sum1 := timedSumFunc(Sum1)
	// sum2 := timedSumFunc(Sum2)
	// fmt.Printf("%d, %d\n", sum1(1, 10000000), sum2(1, 10000000))

	s := strings.NewReader("hello world")
	var b = make([]byte,s.Len())
	s.Read(b)
	fmt.Println(string(b))
	var btA = []string{"ni","是","jia"}
	var btB = make([]string,5)
	copy(btB,btA)
	btA = []string{"nin","jia"}
	copy(btB,btA)
	fmt.Println(btB)
	s.Reset("你好ya")
	bt , _,_ := s.ReadRune()
	fmt.Println(string(bt))
	bt , _,_ = s.ReadRune()
	fmt.Println(string(bt))
	bufS := bufio.NewReader(s)
	ss, err := bufS.ReadBytes('y')
	fmt.Println(bufS.ReadString('a'))
	fmt.Println(string(ss),err)

	copy(b,"aaa")

	//r:= io.NewSectionReader(reader, 1, 4)
	//bufS.WriteTo(s)

	var te = []string{"tes","name"}
	Slice(te)
	fmt.Println(te)

	p:=&sync.Pool{
		New: func() interface{}{
			return 0
		},
	}
	p.Put("jiangzhou")
	p.Put(123456)
	fmt.Println(p.Get())
	fmt.Println(p.Get())
	fmt.Println(p.Get())

}

func Slice(s []string) {
	s[1] = "aaa"
	return
}
