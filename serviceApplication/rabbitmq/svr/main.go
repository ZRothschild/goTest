package main

import (
	"log"
	"os"

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

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	if len(os.Args) < 2 {
		log.Printf("Usage: %s [binding_key]...", os.Args[0])
		os.Exit(0)
	}
	for _, s := range os.Args[1:] {
		log.Printf("Binding queue %s to exchange %s with routing key %s",
			q.Name, "logs_topic", s)
		err = ch.QueueBind(
			q.Name,       // queue name
			s,            // routing key
			"logs_topic", // exchange
			false,
			nil)
		failOnError(err, "Failed to bind a queue")
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for logs. To exit press CTRL+C")
	<-forever
}

//// 订阅模式 消费者消费随机名称的绑定交换机的队列数据，direct 使用 key 来完全匹配消费者,消费同一个
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
//	q, err := ch.QueueDeclare(
//		"",    // name
//		false, // durable
//		false, // delete when unused
//		true,  // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//	failOnError(err, "Failed to declare a queue")
//
//	// 分别启动两个消费者 key 是 dev 和 test
//	for i := 0; i < 2; i++ {
//		if i == 0 {
//			err = ch.QueueBind(
//				q.Name,        // queue name
//				"dev",         // routing key
//				"logs_direct", // exchange
//				false,
//				nil,
//			)
//		} else {
//			err = ch.QueueBind(
//				q.Name,        // queue name
//				"pre",         // routing key
//				"logs_direct", // exchange
//				false,
//				nil,
//			)
//		}
//		failOnError(err, "Failed to bind a queue")
//	}
//
//	msgs, err := ch.Consume(
//		q.Name, // queue
//		"",     // consumer
//		true,   // auto-ack
//		false,  // exclusive
//		false,  // no-local
//		false,  // no-wait
//		nil,    // args
//	)
//	failOnError(err, "Failed to register a consumer")
//
//	var forever chan struct{}
//
//	go func() {
//		for d := range msgs {
//			//fmt.Printf("======== %#+v", d)
//			log.Printf("Received a message: %s\n", d.Body)
//		}
//	}()
//
//	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
//}

//// 订阅模式 消费者消费随机名称的绑定交换机的队列数据，fanout
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
//		"logsExchange", // name
//		"fanout",       // type
//		true,           // durable
//		false,          // auto-deleted
//		false,          // internal
//		false,          // no-wait
//		nil,            // arguments
//	)
//	failOnError(err, "Failed to declare an exchange")
//
//	q, err := ch.QueueDeclare(
//		"",    // name
//		false, // durable
//		false, // delete when unused
//		true,  // exclusive
//		false, // no-wait
//		nil,   // arguments
//	)
//	failOnError(err, "Failed to declare a queue")
//
//	err = ch.QueueBind(
//		q.Name,         // queue name
//		"",             // routing key
//		"logsExchange", // exchange
//		false,
//		nil,
//	)
//	failOnError(err, "Failed to bind a queue")
//
//	msgs, err := ch.Consume(
//		q.Name, // queue
//		"",     // consumer
//		true,   // auto-ack
//		false,  // exclusive
//		false,  // no-local
//		false,  // no-wait
//		nil,    // args
//	)
//	failOnError(err, "Failed to register a consumer")
//
//	var forever chan struct{}
//
//	go func() {
//		for d := range msgs {
//			//fmt.Printf("======== %#+v", d)
//			log.Printf("Received a message: %s\n", d.Body)
//		}
//	}()
//
//	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
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
//	msgs, err := ch.Consume(
//		q.Name, // queue
//		"",     // consumer
//		true,   // auto-ack
//		false,  // exclusive
//		false,  // no-local
//		false,  // no-wait
//		nil,    // args
//	)
//	failOnError(err, "Failed to register a consumer")
//
//	var forever chan struct{}
//
//	go func() {
//		for d := range msgs {
//			//fmt.Printf("======== %#+v", d)
//			log.Printf("Received a message: %s\n", d.Body)
//		}
//	}()
//
//	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
//	<-forever
//}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
