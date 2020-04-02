package main

import (
	"github.com/ZRothschild/goTest/socket/config"
	socketio "github.com/googollee/go-socket.io"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	server, err := socketio.NewServer(nil)
	config.FailOnError(err, "socketio.NewServer")

	client, err := config.MongoClient()
	config.FailOnError(err, "MongoClient")
	collection := client.Database("testing")

	db, err := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
	defer db.Close()
	config.FailOnError(err, "gorm.Open")

	rabbitMq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer rabbitMq.Close()
	config.FailOnError(err, "amqp.Dial")

	config.Sockeio(server, collection, rabbitMq)

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	http.Handle("/dataList", http.HandlerFunc(list))
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func list(w http.ResponseWriter, r *http.Request) {

}
