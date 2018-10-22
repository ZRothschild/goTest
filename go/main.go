package main

import (
	"log"
	"net/http"
	"fmt"
	"io/ioutil"
	"bytes"
	"mime/multipart"
	"os"
	"io"
)
func main() {
	db := database{"shoes":50,"socks":5}
	mux := http.NewServeMux()
	mux.Handle("/list",http.HandlerFunc(db.list))
	mux.Handle("/price",http.HandlerFunc(db.price))
	log.Fatal(http.ListenAndServe(":8000",mux))
}

//金额
type dollars float32

func (d dollars) Strings() string {
	return fmt.Sprintf("$%.2f",	d)
}
//数据结构
type database map[string]dollars
//list
func (db database) list(w http.ResponseWriter, req *http.Request){
	for item,price := range db {
		fmt.Fprintf(w,"%s: %s\n",item,price)
	}

	a,_ := ioutil.ReadAll(req.Body)
	fmt.Fprintf(w,"test %s\n",string(a))
}
//price
func (db database) price(w http.ResponseWriter, req *http.Request){
	item := req.URL.Query().Get("item")
	price,ok :=	db[item]
	if !ok {
		w.WriteHeader(http.StatusNotFound)//404
		fmt.Fprintf(w,"no such item:%q\n",item)
		return
	}
	fmt.Fprintf(w,"%s\n",price)
}


