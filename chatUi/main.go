package main

import (
	"log"
	"net/http"
	"os"
	// "encoding/json"
	// "github.com/ZRothschild/goTest/chatUi/lib"
	// "github.com/ZRothschild/goTest/socket/config"
	// socketio "github.com/googollee/go-socket.io"
	// "github.com/streadway/amqp"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo"
	// "time"
	// "context"
)

func main() {
	serveMux := http.NewServeMux()
	path, _ := os.Getwd()
	hand := http.StripPrefix("/", http.FileServer(http.Dir(path)))
	serveMux.Handle("/", hand)

	// socketServer, err := socketio.NewServer(nil)
	// lib.Log(err, "socketio.NewServer")
	//
	// go socketServer.Serve()
	// defer socketServer.Close()
	// http.Handle("/socket.io/", socketServer)
	//
	// //socket 链接
	// socketServer.OnConnect("/chat", func(s socketio.Conn) error {
	// 	s.Join("Broadcast")
	// 	return nil
	// })
	//
	// //登录事件
	// socketServer.OnEvent("/chat", "login", func(s socketio.Conn, name string) string {
	// 	var result struct {
	// 		Key   string
	// 		Value string
	// 	}
	//
	// 	filter := bson.E{"name", name}
	//
	// 	mongoDb, err := config.MongoClient()
	// 	lib.Log(err, "config.MongoClient")
	//
	// 	numbersCollection := mongoDb.Collection("numbers")
	//
	// 	ctx := context.Background()
	// 	if err := numbersCollection.FindOne(ctx, filter).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
	// 		return err.Error()
	// 	}
	// 	if result.Value == name {
	// 		return "已经登录"
	// 	}
	//
	// 	if _, err := numbersCollection.InsertOne(ctx, filter); err != nil {
	// 		return err.Error()
	// 	}
	// 	return name
	// })
	// //聊天发送信息
	// socketServer.OnEvent("/chat", "msg", func(s socketio.Conn, data map[string]string) {
	//
	// 	dataMsg := bson.M{"name": data["name"], "msg": data["msg"], "add_time": time.Now().Unix()}
	//
	// 	mongoDb, err := config.MongoClient()
	// 	lib.Log(err, "chat config.MongoClient")
	//
	// 	dataMsgCollection := mongoDb.Collection("numbers")
	//
	// 	ctx := context.Background()
	// 	if _, err := dataMsgCollection.InsertOne(ctx, dataMsg); err != nil {
	// 		lib.Log(err, "dataMsgCollection.InsertOne")
	// 	}
	//
	// 	rabbitMq, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// 	defer rabbitMq.Close()
	// 	lib.Log(err, "amqp.Dial")
	//
	// 	ch, err := rabbitMq.Channel()
	// 	lib.Log(err, "rabbitMq.Channel")
	//
	// 	defer ch.Close()
	// 	q, err := ch.QueueDeclare(
	// 		"chatMsgQueue1", // name
	// 		true,            // durable
	// 		false,           // delete when unused
	// 		false,           // exclusive
	// 		false,           // no-wait
	// 		nil,             // arguments
	// 	)
	//
	// 	lib.Log(err, "rabbitMq.Channel")
	//
	// 	msg, err := json.Marshal(data)
	// 	lib.Log(err, "json.Marshal")
	//
	// 	//如果消费端没有 队列go rabbitmq go rabbitmq 绑定交换机 就不要设置交换机名称
	// 	//ch.QueueBind 的第二个参数 将于 ch.Publish 匹配
	// 	err = ch.Publish(
	// 		"",     // exchange
	// 		q.Name, // routing key
	// 		false,  // mandatory
	// 		false,  // immediate
	// 		amqp.Publishing{
	// 			ContentType: "application/json",
	// 			Body:        msg,
	// 		})
	// 	lib.Log(err, "ch.Publish")
	//
	// 	socketServer.BroadcastToRoom("/chat", "Broadcast", "reply", data)
	// })
	//
	// socketServer.OnError("/", func(s socketio.Conn, err error) {
	// 	lib.Log(err, "socketServer.OnError")
	// })
	//
	// socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
	// 	// fmt.Printf("disconect %s\n",reason)
	// })

	server := &http.Server{Addr: "localhost:8080", Handler: serveMux}
	log.Fatal(server.ListenAndServe())
}
