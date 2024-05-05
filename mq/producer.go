package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"encoding/json"
	"github.com/streadway/amqp"
)

func SendUpdateAwardCountMessage(strategyID int64, awardID int) error {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer ch.Close()
	message := map[string]interface{}{
		"strategyID": strategyID,
		"awardID":    awardID,
	}
	marshal, err := json.Marshal(message)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	headers := amqp.Table{
		"x-delay": 3000, // 延时3秒
	}
	err = ch.Publish(
		"delayed-exchange",            // exchange name
		constant.DelayQueueRoutingKey, // routing key
		false,
		false,
		amqp.Publishing{
			Headers:     headers,
			ContentType: "text/plain",
			Body:        marshal,
		},
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	log.Infof("发送消息成功: %v", string(marshal))
	return nil
}
