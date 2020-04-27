package main

import (
	"encoding/json"
	"github.com/ZRothschild/goTest/chatUi/lib"
	"github.com/ZRothschild/goTest/socket/config"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")

	lib.Log(err, "Failed to connect to RabbitMQ")

	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	lib.Log(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"chatMsgQueue1", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	lib.Log(err, "ch.QueueDeclare")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		lib.Log(err, "ch.Qos")
	}

	msgs, err := ch.Consume(
		q.Name,            // queue
		"chatMsgConsume1", // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)

	forever := make(chan bool)

	MySqlDb, err := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		lib.Log(err, "gorm.Open")
	}
	MySqlDb = MySqlDb.LogMode(true) //打印mysql 日子

	defer MySqlDb.Close()

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var test config.Test

			if err = json.Unmarshal(d.Body, &test); err != nil {
				lib.Log(err, "json.Unmarshal")
			}
			if err := MySqlDb.Debug().Table("tests").Create(&test).Error; err != nil {
				lib.Log(err, "插入失败")
			}

			if err = d.Ack(false); err != nil {
				lib.Log(err, "d.Ack")
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
