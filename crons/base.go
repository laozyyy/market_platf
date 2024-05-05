package crons

import (
	"big_market/common/log"
	"big_market/mq"
	"github.com/robfig/cron/v3"
)

func AddCron() {
	c := cron.New(cron.WithSeconds())

	// 每五秒消费一次
	_, err := c.AddFunc("*/5 * * * * *", mq.ConsumeUpdateAwardCountMessage)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	go c.Start()
	defer c.Stop()
}
