package model

import "time"

type RaffleActivityOrder struct {
	//ID            uint64    // 自增ID
	UserID        string    // 用户ID
	SKU           int64     // 商品sku
	ActivityID    int64     // 活动ID
	ActivityName  string    // 活动名称
	StrategyID    int64     // 抽奖策略ID
	OrderID       string    // 订单ID
	OrderTime     time.Time // 下单时间
	TotalCount    int       // 总次数
	DayCount      int       // 日次数
	MonthCount    int       // 月次数
	State         string    // 订单状态（complete）
	OutBusinessNo string    // 业务仿重ID - 外部透传的，确保幂等
	CreateTime    time.Time // 创建时间
	UpdateTime    time.Time // 更新时间
}
