package main

import (
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"fmt"
	"log"
)

func main()  {
	url := "localhost:27017"
	session, err := mgo.Dial(url)
	if err != nil {
		fmt.Printf("err %s\n",err)
	}
	defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("test").C("people")
	err = c.Insert(&Person{"superWang", "13478808311"},
		&Person{"David", "15040268074"})
	if err != nil {
		log.Fatal(err)
	}
	result := Person{}
	err = c.Find(bson.M{"name": "superWang"}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Name:", result.Name)
	fmt.Println("Phone:", result.Phone)
}

type Person struct {
	Name  string
	Phone string
}