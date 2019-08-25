package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Data struct {
	Email string `json:"email"`
	Password string `json:"password"`
	Captcha string `json:"captcha"`
	VerifyString string `json:"verify_string"`
}

func main() {
	db := database{"shoes": 50, "socks": 5}
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/list", db.list)
	serveMux.HandleFunc("/price", db.price)
	serveMux.HandleFunc("/test", db.test)
	hand := http.StripPrefix("/tmp", http.FileServer(http.Dir("F:\\goWork\\goTest\\http")))
	serveMux.Handle("/tmp", hand)
	server := &http.Server{Addr: "localhost:8080", Handler: serveMux}
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

//测试json
func (db database) test (w http.ResponseWriter, req *http.Request) {
	ioReadCloser := req.Body
	fmt.Printf("ioReadCloser %v\n",ioReadCloser)

	tetIo := ioReadCloser
	//读取了 ioReadCloser 就会为空 可以看 req 变化就可知
	reqBody,err := ioutil.ReadAll(ioReadCloser)
	if err != nil {
		fmt.Errorf("%s\n",err)
	}
	var reqData Data
	err = json.Unmarshal(reqBody,&reqData)
	if err != nil {
		fmt.Printf("reqData => %s\n",err)
	}
	fmt.Printf("%v\n",reqData)

	var mapData map[string]string
	err = json.Unmarshal(reqBody,&mapData)
	if err != nil {
		fmt.Errorf("mapData %s\n",err)
	}
	fmt.Printf("%v\n",mapData)

	var reqDataDe Data
	//是取地址改变了  ioReadCloser buff 已经被读取过了
	json.NewDecoder(tetIo).Decode(&reqDataDe)
	fmt.Printf("reqDataDe => %v\n",reqDataDe)

	reqStr := string(reqBody)
	fmt.Printf("%s\n",reqStr)
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
