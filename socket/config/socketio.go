package config

import (
	"context"
	"fmt"
	socketio "github.com/googollee/go-socket.io"
	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

func Sockeio(server *socketio.Server, mongoCollection *mongo.Database, rabbitMq *amqp.Connection) {
	ctx := context.Background()

	server.OnConnect("/chat", func(s socketio.Conn) error {
		s.SetContext("")
		fmt.Println("connected:", s.ID())
		s.Join("Broadcast")
		return nil
	})

	server.OnEvent("/chat", "login", func(s socketio.Conn, name string) string {
		filter := bson.E{"name", name}
		var result struct {
			Key   string
			Value string
		}
		numbersCollection := mongoCollection.Collection("numbers")
		err := numbersCollection.FindOne(ctx, filter).Decode(&result)
		if err != nil && err != mongo.ErrNoDocuments {
			return err.Error()
		}
		if result.Value == name {
			return "已经登录"
		}
		_, err = numbersCollection.InsertOne(ctx, filter)
		if err != nil {
			return err.Error()
		}
		return name
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, data map[string]string) {
		dataMsg := bson.M{"name": data["name"], "msg": data["msg"], "time": time.Now().Unix()}
		dataMsgCollection := mongoCollection.Collection("numbers")
		_, err := dataMsgCollection.InsertOne(ctx, dataMsg)
		if err != nil {
			fmt.Printf("dataMsgCollection.InsertOne => %s", err)
		}
		//rabbitmq 异步写入数据
		err = SendMsg(rabbitMq, dataMsg)
		if err != nil {
			fmt.Printf("SendMsg => %s", err)
		}

		server.BroadcastToRoom("/chat", "Broadcast", "reply", data)
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})
}
