package main

import (
	"encoding/json"
	"github.com/ZRothschild/goTest/socket/config"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	config.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	msgs := config.ConsumeMsg(conn)
	forever := make(chan bool)
	db, _ := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true) //打印mysql 日子
	defer db.Close()
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var test config.Test
			err = json.Unmarshal(d.Body, &test)
			if err != nil {
				config.FailOnError(err, "json.Unmarshal")
			}
			if err := db.Debug().Table("tests").Create(&test).Error; err != nil {
				config.FailOnError(err, "插入失败")
			}
			err = d.Ack(false)
			config.FailOnError(err, "d.Ack")
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
