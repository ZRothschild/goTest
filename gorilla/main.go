package main

import (
	"encoding/json"
	"fmt"
)

type Message struct {
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
	Age       int    `json:"age,omitempty"`
}

func main() {
	msg := Message{Sender: "test", Recipient: "aa", Age: 0}
	jsonStr, _ := json.Marshal(msg)
	fmt.Printf("%s", jsonStr)
}
