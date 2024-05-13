package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/service/reposity"
	"encoding/json"
)

// ConsumeUpdateAwardCountMessage 拉模式
func ConsumeUpdateAwardCountMessage() {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	defer ch.Close()
	// 每次拉取一条消息
	msg, ok, err := ch.Get(constant.DelayQueueName, false)
	if ok {
		log.Infof("UpdateAwardCount消费者收到消息: %s", string(msg.Body))
		message := make(map[string]interface{})
		err = json.Unmarshal(msg.Body, &message)
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
		// 处理消息的逻辑
		strategyID := message["strategyID"].(float64)
		awardID := message["awardID"].(float64)
		err = reposity.UpdateStrategyAwardCount(int64(strategyID), int(awardID))
		if err != nil {
			log.Errorf("消费失败, strategyID: %d, awardID: %d", strategyID, awardID)
			return
		}
		// 手动确认消息
		err = msg.Ack(false)
		if err != nil {
			log.Errorf("err: %v", err)
			return
		}
		log.Infof("消费成功")
	} else {
		log.Infof("队列中无消息")
	}
}

// ConsumeUpdateSkuCountMessage 推模式
func ConsumeUpdateSkuCountMessage() {
	ch, err := Conn.Channel()
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	defer ch.Close()
	// 消费消息
	consumerTag := "my_consumer"

	msgs, err := ch.Consume(
		constant.UpdateSkuCountQueue, // 队列名称
		consumerTag,                  // 消费者标签
		false,                        // 是否自动确认消息
		false,                        // 是否独占消费者（仅限于本连接）
		false,                        // 是否阻塞等待服务器确认
		false,                        // 是否使用内部排他队列
		nil,                          // 其他参数
	)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}

	// 启动消费者协程
	go func() {
		log.Infof("UpdateSkuCount消费者协程启动")
		for msg := range msgs {
			log.Infof("UpdateSkuCount消费者收到消息: %s", string(msg.Body))
			// 手动确认消息已被消费
			msg.Ack(false)
		}
	}()
}
