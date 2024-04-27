package model

import "time"

type Award struct {
	ID          int       `json:"id"`
	AwardID     int       `json:"awardId"`
	AwardKey    string    `json:"awardKey"`
	AwardConfig string    `json:"awardConfig"`
	AwardDesc   string    `json:"awardDesc"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}
