package main

import (
	"encoding/json"
	"github.com/ZRothschild/goTest/chatUi/lib"
	"github.com/ZRothschild/goTest/socket/config"
	socketio "github.com/googollee/go-socket.io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
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

var (
	userTable string = "test"
)

type User struct {
	Email  string `json:"email"`
	Avatar string `json:"avatar"`
}

func main() {
	serveMux := http.NewServeMux()
	// 用户注册
	serveMux.HandleFunc("/login", func(writer http.ResponseWriter, request *http.Request) {
		var user User
		if err := json.NewDecoder(request.Body).Decode(&user); err != nil {
			lib.Log(err, "json.NewDecoder(request.Body)")
			return
		}
		var imgArr = []string{"/img/2.jpeg", "/img/3.jpeg", "./img/1.png"}
		rand.Seed(time.Now().UnixNano())
		user.Avatar = imgArr[rand.Intn(3)]

		if err := config.MySqlDb().Debug().Table(userTable).Create(&userTable).Error; err != nil {
			config.FailOnError(err, "插入失败")
		}

		writer.Header().Set("content-type", "application/json;charset=utf-8")
		if err := json.NewEncoder(writer).Encode(user); err != nil {
			lib.Log(err, "json.NewEncoder(writer)")
			return
		}
	})
	path, _ := os.Getwd()
	hand := http.StripPrefix("/", http.FileServer(http.Dir(path)))
	serveMux.Handle("/", hand)

	socketServer, err := socketio.NewServer(nil)
	lib.Log(err, "socketio.NewServer")

	defer socketServer.Close()
	// socket 链接
	socketServer.OnConnect("/chat", func(s socketio.Conn) error {
		s.Join("Broadcast")
		return nil
	})

	// 登录事件
	socketServer.OnEvent("/chat", "login", func(s socketio.Conn, name string) string {
		return name
	})
	// 聊天发送信息
	socketServer.OnEvent("/chat", "msg", func(s socketio.Conn, data map[string]string) {
		socketServer.BroadcastToRoom("/chat", "Broadcast", "reply", data)
	})

	socketServer.OnError("/", func(s socketio.Conn, err error) {
		lib.Log(err, "socketServer.OnError")
	})

	socketServer.OnDisconnect("/", func(s socketio.Conn, reason string) {
		// fmt.Printf("disconnect %s\n",reason)
	})

	go socketServer.Serve()
	http.Handle("/socket.io/", socketServer)
	server := &http.Server{Addr: "localhost:8080", Handler: serveMux}
	log.Fatal(server.ListenAndServe())
}
