package model

import "time"

type RaffleActivity struct {
	ID            uint64    `json:"id"`              // 自增ID
	ActivityID    int64     `json:"activity_id"`     // 活动ID
	ActivityName  string    `json:"activity_name"`   // 活动名称
	ActivityDesc  string    `json:"activity_desc"`   // 活动描述
	BeginDateTime time.Time `json:"begin_date_time"` // 开始时间
	EndDateTime   time.Time `json:"end_date_time"`   // 结束时间
	StrategyID    int64     `json:"strategy_id"`     // 抽奖策略ID
	State         string    `json:"state"`           // 活动状态
	CreateTime    time.Time `json:"create_time"`     // 创建时间
	UpdateTime    time.Time `json:"update_time"`     // 更新时间
}
