package main

import (
	"fmt"
	"reflect"
	"strconv"
)

func main() {

	// var i int = 2
	// //Elem 返回值
	// inte := reflect.ValueOf(&i).Elem()
	// fmt.Printf("%v,%T\n",inte,inte)     //2，ptr

	myType := &MyType{22, "my name is liLei"}
	// reflect.ValueOf 返回地址
	// reflect.ValueOf(&myType).Elem() 返回值  &{22 my name is liLei}
	mtV := reflect.ValueOf(&myType).Elem()
	// 触发 String 方法   0xc000052400--name:my name is liLei i:22  =====  reflect.Value
	fmt.Printf("%v  =====  %T\n", mtV, mtV)
	// 因为打印是有返回值所有去第零个返回值
	fmt.Println("Before:", mtV.MethodByName("String").Call(nil)[0])

	params := make([]reflect.Value, 1)
	params[0] = reflect.ValueOf(18)
	mtV.MethodByName("SetI").Call(params)
	params[0] = reflect.ValueOf("reflection test")
	mtV.MethodByName("SetName").Call(params)
	fmt.Println("After:", mtV.MethodByName("String").Call(nil)[0])
	fmt.Println("After:", mtV.Method(2).Call(nil)[0])
}

type MyType struct {
	i    int
	name string
}

func (mt *MyType) SetI(i int) {
	mt.i = i
}

func (mt *MyType) SetName(name string) {
	mt.name = name
}

func (mt *MyType) String() string {
	return fmt.Sprintf("%p", mt) + "--name:" + mt.name + " i:" + strconv.Itoa(mt.i)
}
