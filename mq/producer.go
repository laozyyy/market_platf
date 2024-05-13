package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
)

func SendMessage(exchange string, topic string, publish amqp.Publishing) error {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	defer ch.Close()
	err = ch.Publish(
		exchange, // exchange name
		topic,    // routing key
		false,
		false,
		publish,
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	log.Infof("发送消息成功: %v", string(publish.Body))
	return nil
}

func SendUpdateAwardCountMessage(strategyID int64, awardID int) error {
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
	publish := amqp.Publishing{
		Headers:     headers,
		ContentType: "text/plain",
		Body:        marshal,
	}
	log.Infof("发送更新奖品库存延迟消息")
	err = SendMessage(constant.DelayExchangeName, constant.UpdateStrategyAwardCountTopic, publish)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

func SendSkuCountZeroMessage(sku int64) error {
	message := map[string]interface{}{
		"sku": sku,
	}
	marshal, err := json.Marshal(message)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	publish := amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
	}
	log.Infof("发送活动sku库存归零消息")
	err = SendMessage(constant.NormalExchangeName, constant.SkuCountZeroTopic, publish)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

// SendUpdateSkuCountMessage 延时发送
func SendUpdateSkuCountMessage(sku int64) error {
	message := map[string]interface{}{
		"sku": sku,
	}
	marshal, err := json.Marshal(message)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	headers := amqp.Table{
		"x-delay": 3000, // 延时3秒
	}
	publish := amqp.Publishing{
		ContentType: "text/plain",
		Body:        marshal,
		Headers:     headers,
	}
	log.Infof("发送更新活动sku库存延迟消息")
	err = SendMessage(constant.DelayExchangeName, constant.UpdateSkuCountTopic, publish)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
