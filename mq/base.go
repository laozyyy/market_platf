package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/conf"
	"github.com/streadway/amqp"
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
	// 声明正常队列
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
	// 绑定正常队列到延时交换机
	err = ch.QueueBind(
		constant.DelayQueueName,       // queue name
		constant.DelayQueueRoutingKey, // routing key
		constant.DelayExchangeName,    // exchange name
		false,                         // no-wait
		nil,                           // arguments
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
}
