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

func SocketIo(server *socketio.Server, mongoDb *mongo.Database, rabbitMq *amqp.Connection) {
	ctx := context.Background()

	server.OnConnect("/chat", func(s socketio.Conn) error {
		s.Join("Broadcast")
		return nil
	})

	server.OnEvent("/chat", "login", func(s socketio.Conn, name string) string {
		filter := bson.E{"name", name}
		var result struct {
			Key   string
			Value string
		}
		numbersCollection := mongoDb.Collection("numbers")
		if err := numbersCollection.FindOne(ctx, filter).Decode(&result); err != nil && err != mongo.ErrNoDocuments {
			return err.Error()
		}
		if result.Value == name {
			return "已经登录"
		}
		_, err := numbersCollection.InsertOne(ctx, filter)
		if err != nil {
			return err.Error()
		}
		return name
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, data map[string]string) {
		dataMsg := bson.M{"name": data["name"], "msg": data["msg"], "add_time": time.Now().Unix()}
		dataMsgCollection := mongoDb.Collection("numbers")
		_, err := dataMsgCollection.InsertOne(ctx, dataMsg)
		if err != nil {
			FailOnError(err, "dataMsgCollection.InsertOne")
		}

		if err = SendMsg(rabbitMq, dataMsg); err != nil { //rabbitmq 异步写入数据
			FailOnError(err, "rabbitMq SendMsg")
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
