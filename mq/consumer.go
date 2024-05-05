package mq

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/service/reposity"
	"encoding/json"
)

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
		log.Infof("消费者收到消息: %s", string(msg.Body))
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
