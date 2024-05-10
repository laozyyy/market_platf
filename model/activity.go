package model

import "time"

type SkuRechargeEntity struct {
	UserID        string `json:"userID"`
	Sku           int64  `json:"sku"`
	OutBusinessNo string `json:"outBusinessNo"`
}

type ActivityOrderEntity struct {
	UserID        string // 用户ID
	Sku           int64
	ActivityID    int64     // 活动ID
	ActivityName  string    // 活动名称
	StrategyID    int64     // 抽奖策略ID
	OrderID       string    // 订单ID
	OrderTime     time.Time // 下单时间
	TotalCount    int       // 总次数
	DayCount      int       // 日次数
	MonthCount    int       // 月次数
	State         string    // 订单状态
	OutBusinessNo string
}

type CreateOrderAggregate struct {
	UserID              string
	ActivityID          int64
	TotalCount          int
	DayCount            int
	MonthCount          int
	ActivityOrderEntity ActivityOrderEntity
}
