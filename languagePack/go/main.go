package main

import (
	"fmt"
	"io/ioutil"
	_ "mime/multipart"
	"net/http"
	"reflect"
	"unsafe"
)

type W struct {
	a byte
	b int32
	c int64
}

type T struct {
	t1 byte
	t2 int16
	t3 int64
	t4 []int
	// t5 bool
}

func main() {
	var w *W
	var w2 W
	//取值类型的大小 指针类型 是8 或者4
	fmt.Println(unsafe.Sizeof(w)) // 8
	fmt.Printf("%v\n", w)

	fmt.Println(unsafe.Sizeof(*w)) //结构体是16
	fmt.Println(unsafe.Sizeof(&w)) //8

	//对齐位为8 所以w sizeof 是16 才能容下 w 如果 w.c 改成int32 就是4 对齐位
	fmt.Println(unsafe.Alignof(w.a)) // type byte  1
	fmt.Println(unsafe.Alignof(w.b)) // type int32 4
	fmt.Println(unsafe.Alignof(w.c)) // type int64 8

	fmt.Println(unsafe.Alignof(w))                //8
	fmt.Println(unsafe.Alignof(w2))               //8
	fmt.Println(reflect.TypeOf(w).Elem().Align()) //8

	fmt.Println("----------unsafe.Pointer---------")
	t := &T{1, 2, 3, []int{1, 2}}
	ptr := unsafe.Pointer(t)
	t1 := (*byte)(ptr)
	fmt.Println(*t1)
	t2 := (*int32)(unsafe.Pointer(uintptr(ptr) + unsafe.Offsetof(t.t2)))

	fmt.Printf("t2 === %d\n", t2)
	*t2 = 99
	fmt.Println(t)
	fmt.Println("===============")
	fmt.Println((*int32)(unsafe.Pointer(unsafe.Offsetof(t.t2))))
	fmt.Println("===============")
	t3 := (*int64)(unsafe.Pointer(uintptr(ptr) + unsafe.Offsetof(t.t3)))
	fmt.Printf("t3 === %d\n", t3)
	*t3 = 123
	fmt.Println(t)
	fmt.Println("##################")
	fmt.Println(unsafe.Alignof(t.t4))
	fmt.Println((*int32)(unsafe.Pointer(unsafe.Offsetof(t.t1))))
	fmt.Println((*int32)(unsafe.Pointer(unsafe.Offsetof(t.t2))))
	fmt.Println((*int32)(unsafe.Pointer(unsafe.Offsetof(t.t3))))
	fmt.Println((*int32)(unsafe.Pointer(unsafe.Offsetof(t.t4))))
	fmt.Println(unsafe.Offsetof(t.t4))
	// fmt.Println( (*int32)(unsafe.Pointer(unsafe.Offsetof(t.t5))))
	fmt.Println("##################")

	/**
		var x struct {
		a bool
		b int16
		c []int
	}

	Sizeof(x)   = 32  Alignof(x)   = 8
	Sizeof(x.a) = 1   Alignof(x.a) = 1 Offsetof(x.a) = 0
	Sizeof(x.b) = 2   Alignof(x.b) = 2 Offsetof(x.b) = 2
	Sizeof(x.c) = 24  Alignof(x.c) = 8 Offsetof(x.c) = 8
	*/

	// db := database{"shoes": 50, "socks": 5}
	// mux := http.NewServeMux()
	// mux.Handle("/list", http.HandlerFunc(db.list))
	// mux.Handle("/price", http.HandlerFunc(db.price))
	// log.Fatal(http.ListenAndServe(":8000", mux))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f", d)
}

//数据结构
type database map[string]dollars

//list
func (db database) list(w http.ResponseWriter, req *http.Request) {
	for item, price := range db {
		fmt.Fprintf(w, "%s: %s\n", item, price)
	}

	a, _ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(w, "test %s\n", string(a))
}

//price
func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price, ok := db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound) //404
		fmt.Fprintf(w, "no such item:%q\n", item)
		return
	}
	fmt.Fprintf(w, "%s\n", price)
}
