package model

import "time"

type RaffleActivityCount struct {
	ID              uint64    `json:"id"`                // 自增ID
	ActivityCountID int64     `json:"activity_count_id"` // 活动次数编号
	TotalCount      int       `json:"total_count"`       // 总次数
	DayCount        int       `json:"day_count"`         // 日次数
	MonthCount      int       `json:"month_count"`       // 月次数
	CreateTime      time.Time `json:"create_time"`       // 创建时间
	UpdateTime      time.Time `json:"update_time"`       // 更新时间
}
