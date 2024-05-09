package model

import "time"

type RaffleActivitySku struct {
	ID                uint64    `json:"id"`                  // 自增ID
	SKU               int64     `json:"sku"`                 // 商品sku - 把每一个组合当做一个商品
	ActivityID        int64     `json:"activity_id"`         // 活动ID
	ActivityCountID   int64     `json:"activity_count_id"`   // 活动个人参与次数ID
	StockCount        int       `json:"stock_count"`         // 商品库存
	StockCountSurplus int       `json:"stock_count_surplus"` // 剩余库存
	CreateTime        time.Time `json:"create_time"`         // 创建时间
	UpdateTime        time.Time `json:"update_time"`         // 更新时间
}
