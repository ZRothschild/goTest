package main

import (
	"github.com/ZRothschild/goTest/socket/config"
	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func main() {
	server, err := socketio.NewServer(nil)
	config.FailOnError(err, "socketio.NewServer")

	mongoDb, err := config.MongoClient()
	config.FailOnError(err, "MongoClient")

	rabbitMq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	defer rabbitMq.Close()
	config.FailOnError(err, "amqp.Dial")

	config.SocketIo(server, mongoDb, rabbitMq)

	go server.Serve()
	defer server.Close()
	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	http.Handle("/dataList", http.HandlerFunc(config.List))
	log.Println("Serving at localhost:8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
