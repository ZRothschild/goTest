package web

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"net/http"
	"strings"
)

//全局信息
var datas Datas
var users map[*websocket.Conn]string

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./view/test.html")
}

func webSocket(ws *websocket.Conn) {
	var userMsg UserMsg
	var data string
	for {
		//判断是否重复连接
		if _, ok := users[ws]; !ok {
			users[ws] = "匿名"
		}
		userMsgLen := len(datas.UserMsgs)
		fmt.Println("UserMsg", userMsgLen, "users长度：", len(users))
		//有消息时，全部分发送数据
		if userMsgLen > 0 {
			b, errMarshal := json.Marshal(datas)
			if errMarshal != nil {
				fmt.Println("全局消息内容异常...")
				break
			}
			for key, _ := range users {
				errMarshal = websocket.Message.Send(key, string(b))
				if errMarshal != nil {
					//移除出错的链接
					delete(users, key)
					fmt.Println("发送出错...")
					break
				}
			}
			datas.UserMsgs = make([]UserMsg, 0)
		}
		fmt.Println("开始解析数据...")
		err := websocket.Message.Receive(ws, &data)
		fmt.Println("data：", data)
		if err != nil {
			//移除出错的链接
			delete(users, ws)
			fmt.Println("接收出错...")
			break
		}
		data = strings.Replace(data, "\n", "", 0)
		err = json.Unmarshal([]byte(data), &userMsg)
		if err != nil {
			fmt.Println("解析数据异常...")
			break
		}
		fmt.Println("请求数据类型：", userMsg.DataType)
		switch userMsg.DataType {
		case "send":
			//赋值对应的昵称到ws
			if _, ok := users[ws]; ok {
				users[ws] = userMsg.UserName
				//清除连接人昵称信息
				datas.UserDatas = make([]UserData, 0)
				//重新加载当前在线连接人
				for _, item := range users {
					userData := UserData{UserName: item}
					datas.UserDatas = append(datas.UserDatas, userData)
				}
			}
			datas.UserMsgs = append(datas.UserMsgs, userMsg)
		}
	}
}

type UserMsg struct {
	UserName string
	Msg      string
	DataType string
}

type UserData struct {
	UserName string
}

type Datas struct {
	UserMsgs  []UserMsg
	UserDatas []UserData
}
