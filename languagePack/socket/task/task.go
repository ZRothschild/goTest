package main

import (
	"encoding/json"
	"github.com/ZRothschild/goTest/socket/config"
	"github.com/streadway/amqp"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	config.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	defer ch.Close()

	config.FailOnError(err, "Failed to open a channel")
	msgs := config.ConsumeMsg(ch)

	forever := make(chan bool)
	defer config.MySqlDb().Close()

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			var test config.Test

			if err = json.Unmarshal(d.Body, &test); err != nil {
				config.FailOnError(err, "json.Unmarshal")
			}
			if err := config.MySqlDb().Debug().Table("tests").Create(&test).Error; err != nil {
				config.FailOnError(err, "插入失败")
			}

			if err = d.Ack(false); err != nil {
				config.FailOnError(err, "d.Ack")
			}
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
