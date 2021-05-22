//把数据渲染到浏览器
package main

import (
	"fmt"
	"text/template"
	"net/http"
)

//定义全局的模板变量
var mytemplate *template.Template

type User struct {
	Name string
	Role string
}
func init () {
	var err error
	mytemplate, err = template.ParseFiles("./index.html")
	if err != nil {
		fmt.Println("parse file failed, error:", err)
	}

}

func dealTest1 (w http.ResponseWriter, r *http.Request) {
	u1 := User{
		"admin",
		"管理员",
	}
	err := mytemplate.Execute(w, u1)
	if err != nil{
		fmt.Println("excute failed, error:", err)
	}
}

func fromData (w http.ResponseWriter, r *http.Request) {
	fname := r.URL.Query().Get("fname")
	lname := r.URL.Query().Get("lname")
	fmt.Println(fname,lname)
}

func main () {
	http.HandleFunc("/test1", dealTest1)
	http.HandleFunc("/data", fromData)
	err := http.ListenAndServe("localhost:8080", nil)
	if err != nil {
		fmt.Println("listen server failed, error:", err)
		return
	}
}
