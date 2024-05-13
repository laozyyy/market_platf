package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/conf"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	Conn *amqp.Connection
	url  string
)

func getMQConnection() (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return conn, nil
}

func init() {
	config := conf.LoadConfig()
	url = config.MQ.URL
}

func Init() {
	initExchangeAndQueue()
	ConsumeUpdateSkuCountMessage()
}

func initExchangeAndQueue() {
	var err error
	Conn, err = getMQConnection()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	ch, err := Conn.Channel()
	defer ch.Close()
	// 交换机
	err = ch.ExchangeDeclare(
		constant.DelayExchangeName, // exchange name
		"x-delayed-message",        // exchange type
		true,                       // durable
		false,                      // auto-deleted
		false,                      // internal
		false,                      // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = ch.ExchangeDeclare(
		constant.NormalExchangeName, // exchange name
		"direct",                    // exchange type
		true,                        // durable
		false,                       // auto-deleted
		false,                       // internal
		false,                       // no-wait
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	// 声明延时队列
	_, err = ch.QueueDeclare(
		constant.DelayQueueName, // queue name
		true,                    // durable
		false,                   // delete when unused
		false,                   // exclusive
		false,                   // no-wait
		nil,                     // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	_, err = ch.QueueDeclare(
		constant.UpdateSkuCountQueue, // queue name
		true,                         // durable
		false,                        // delete when unused
		false,                        // exclusive
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	// 绑定2队列到延时交换机
	err = ch.QueueBind(
		constant.DelayQueueName,                // queue name
		constant.UpdateStrategyAwardCountTopic, // routing key
		constant.DelayExchangeName,             // exchange name
		false,                                  // no-wait
		nil,                                    // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = ch.QueueBind(
		constant.UpdateSkuCountQueue, // queue name
		constant.UpdateSkuCountTopic, // routing key
		constant.DelayExchangeName,   // exchange name
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	err = ch.QueueBind(
		constant.UpdateSkuCountQueue, // queue name
		constant.SkuCountZeroTopic,   // routing key
		constant.NormalExchangeName,  // exchange name
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
}
