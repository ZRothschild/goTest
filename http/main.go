package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/list", db.list)
	serveMux.HandleFunc("/price", db.price)
	hand := http.StripPrefix("/tmp", http.FileServer(http.Dir("F:\\goWork\\goTest\\http")))
	serveMux.Handle("/tmp", hand)
	server := &http.Server{Addr: "localhost:80", Handler: serveMux}
	log.Fatal(server.ListenAndServe())
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
		fmt.Fprintf(w, "%s: %f\n", item, price)
	}
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
	fmt.Fprintf(w, "%f\n", price)
}

// mux := http.NewServeMux()
// mux.Handle("/list",http.HandlerFunc(db.list))
// mux.Handle("/list",http.HandlerFunc(db.list))
// log.Fatal(http.ListenAndServe(":8000",mux))

// http.HandleFunc("/list",db.list)
// http.HandleFunc("/price",db.price)
// log.Fatal(http.ListenAndServe("localhost:8000",	nil))

//FileServer 文件服务器使用方法
// http.Handle("/", http.FileServer(http.Dir("F:\\goWork\\goTest\\http")))
// http.Handle("/tmp", http.StripPrefix("/tmp", http.FileServer(http.Dir("F:\\goWork\\goTest\\http"))))
// http.ListenAndServe(":80", nil)
