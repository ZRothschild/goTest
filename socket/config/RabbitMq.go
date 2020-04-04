package config

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

//生产消息
func SendMsg(conn *amqp.Connection, data map[string]interface{}) error {
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chatMsgQueue1", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		return err
	}
	msg, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//如果消费端没有 队列go rabbitmq go rabbitmq 绑定交换机 就不要设置交换机名称
	//ch.QueueBind 的第二个参数 将于 ch.Publish 匹配
	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	return err
}

//消费消息
func ConsumeMsg(conn *amqp.Connection) <-chan amqp.Delivery {
	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chatMsgQueue1", // name
		true,            // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	FailOnError(err, "Failed to set QoS")

	msg, err := ch.Consume(
		q.Name,            // queue
		"chatMsgConsume1", // consumer
		false,             // auto-ack
		false,             // exclusive
		false,             // no-local
		false,             // no-wait
		nil,               // args
	)
	return msg
}
