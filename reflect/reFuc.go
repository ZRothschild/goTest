package main

import (
	"fmt"
	"github.com/iGoogle-ink/gopay"
	"reflect"
)

func main() {

	// 函数值
	// test := test
	// fmt.Printf("%T %v\n",test,test)
	pStr := "你好北京"
	// test(pStr)

	// 判断 test 类型
	fv := reflect.ValueOf(test)
	// fv.Kind() 判断是否为函数，否则报错
	fmt.Println("fv is reflect.Func ?", fv.Kind() == reflect.Func)
	// 生成reflect.Value 组装数据
	intParams := make([]reflect.Value, 0)
	intParams = append(intParams, reflect.ValueOf(pStr))
	// 注意参数个数与函数接收参数保持一致
	// intParams := make([]reflect.Value,1)
	// intParams[0] = intParams,reflect.ValueOf(pStr)
	fvRe := fv.Call(intParams)
	fmt.Printf("没有返回值 %T\n", fvRe) // []reflect.Value

	/**************************************************************************************/

	// 多参数传值
	testMany := reflect.ValueOf(testMany)
	if testMany.Kind() == reflect.Func {
		intMany := make([]reflect.Value, 0)
		intMany = append(intMany, reflect.ValueOf(pStr), reflect.ValueOf(pStr), reflect.ValueOf(pStr))
		testRes := testMany.Call(intMany)
		// Interface 变成接口类型
		fmt.Printf("返回值 %T %T %T\n",
			testRes,
			testRes[0].Interface(),
			testRes[0].Interface().(string))
	}
}

func test(str string) {
	var s string = "单个参数"
	s = s + str
	fmt.Println(s)
}

func testMany(str ...string) string {
	var s string = "多参数"
	// str []string
	fmt.Printf("%T %v\n", str, str)
	for _, v := range str {
		s = s + v
	}
	fmt.Println(s)
	return s
}
