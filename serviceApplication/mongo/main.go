package main

import (
	"fmt"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func main() {
	url := "mongodb://zhao:123456@localhost:27017/go"
	//默认链接类型
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n", err)
	}
	defer session.Close()
	//err = session.DB("go").Login("zhao","123456")
	//fmt.Printf("loin %v\n",err)
	db := session.DB("go").C("test")
	var res map[string]interface{}
	err = db.Find(bson.M{"lang": "go"}).One(&res)
	fmt.Printf("res %v\n", res)
	fmt.Printf("err %v\n", err)
}
