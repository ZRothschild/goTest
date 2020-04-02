package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

type Test struct {
	Id      int64 `gorm:"-;primary_key;AUTO_INCREMENT"`
	Msg     string
	Name    string
	AddTime int64
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)

	db, _ := gorm.Open("mysql", "root:Nm123456.@/test?charset=utf8&parseTime=True&loc=Local")
	db.LogMode(true)
	defer db.Close()
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var test Test
			err = json.Unmarshal(d.Body, &test)
			if err != nil {
				log.Printf("Received a err: %s", err)
			}
			if err := db.Debug().Table("tests").Create(&test).Error; err != nil {
				fmt.Println("插入失败")
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
