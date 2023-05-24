package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/my_vhost")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := bodyFrom(os.Args)
	err = ch.PublishWithContext(ctx,
		"logs_topic",          // exchange
		severityFrom(os.Args), // routing key
		false,                 // mandatory
		false,                 // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	failOnError(err, "Failed to publish a message")

	log.Printf(" [x] Sent %s", body)
}

func bodyFrom(args []string) string {
	var s string
	if (len(args) < 3) || os.Args[2] == "" {
		s = "hello"
	} else {
		s = strings.Join(args[2:], " ")
	}
	return s
}

func severityFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}

//// 订阅模式 生产者直接生产到交换机 direct 使用 key 来完全匹配消费者 消费同一个
//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/my_vhost")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	err = ch.ExchangeDeclare(
//		"logs_direct", // name
//		"direct",      // type
//		true,          // durable
//		false,         // auto-deleted
//		false,         // internal
//		false,         // no-wait
//		nil,           // arguments
//	)
//	failOnError(err, "Failed to declare an exchange")
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	body := "Hello World!"
//	keys := []string{"test", "test", "pre", "dev", "pre", "dev", "test", "dev", "test", "dev"}
//
//	for i := 0; i < 10; i++ {
//		var by string
//		by = body + "#=>" + strconv.Itoa(i) + "#" + keys[i]
//		err = ch.PublishWithContext(ctx,
//			"logs_direct", // exchange
//			keys[i],       // routing key
//			false,         // mandatory
//			false,         // immediate
//			amqp.Publishing{
//				ContentType: "text/plain",
//				Body:        []byte(by),
//			})
//		failOnError(err, "Failed to publish a message")
//		log.Printf(" [x] Sent %s\n", by)
//	}
//}

//// 订阅模式 生产者直接生产到交换机 fanout
//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/my_vhost")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//	if err = ch.ExchangeDeclare(
//		"logsExchange", // name
//		"fanout",       // type
//		true,           // durable
//		false,          // auto-deleted
//		false,          // internal
//		false,          // no-wait
//		nil,            // arguments
//	); err != nil {
//		failOnError(err, "Failed to ExchangeDeclare a message")
//	}
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	body := "Hello World!"
//
//	for i := 0; i < 10; i++ {
//		var by string
//		by = body + "#=>" + strconv.Itoa(i)
//		err = ch.PublishWithContext(ctx,
//			"logsExchange", // exchange
//			"",             // routing key
//			false,          // mandatory
//			false,          // immediate
//			amqp.Publishing{
//				ContentType: "text/plain",
//				Body:        []byte(by),
//			})
//		failOnError(err, "Failed to publish a message")
//		log.Printf(" [x] Sent %s\n", by)
//	}
//}

//// 简单类型没有交换机
//func main() {
//	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/my_vhost")
//	failOnError(err, "Failed to connect to RabbitMQ")
//	defer conn.Close()
//
//	ch, err := conn.Channel()
//	failOnError(err, "Failed to open a channel")
//	defer ch.Close()
//
//
//	q, err := ch.QueueDeclare(
//		"hello", // name
//		false,   // durable
//		false,   // delete when unused
//		false,   // exclusive
//		false,   // no-wait
//		nil,     // arguments
//	)
//	failOnError(err, "Failed to declare a queue")
//
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//
//	body := "Hello World! zhaoq"
//
//	for i := 0; i < 10; i++ {
//		var by string
//		by = body + "#=>" + strconv.Itoa(i)
//		err = ch.PublishWithContext(ctx,
//			"",     // exchange
//			q.Name, // routing key
//			false,  // mandatory
//			false,  // immediate
//			amqp.Publishing{
//				ContentType: "text/plain",
//				Body:        []byte(by),
//			})
//		failOnError(err, "Failed to publish a message")
//		log.Printf(" [x] Sent %s\n", by)
//	}
//}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
