package model

import "time"

type RaffleActivityAccount struct {
	ID                uint64    `json:"id"`                  // 自增ID
	UserID            string    `json:"user_id"`             // 用户ID
	ActivityID        int64     `json:"activity_id"`         // 活动ID
	TotalCount        int       `json:"total_count"`         // 总次数
	TotalCountSurplus int       `json:"total_count_surplus"` // 总次数-剩余
	DayCount          int       `json:"day_count"`           // 日次数
	DayCountSurplus   int       `json:"day_count_surplus"`   // 日次数-剩余
	MonthCount        int       `json:"month_count"`         // 月次数
	MonthCountSurplus int       `json:"month_count_surplus"` // 月次数-剩余
	CreateTime        time.Time `json:"create_time"`         // 创建时间
	UpdateTime        time.Time `json:"update_time"`         // 更新时间
}
