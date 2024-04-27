package model

import "time"

type Strategy struct {
	ID           int       `json:"id"`
	StrategyID   int       `json:"strategyId"`
	StrategyDesc string    `json:"strategyDesc"`
	RuleModels   string    `json:"ruleModels"` //规则模型，rule配置的模型同步到此表，便于使用
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}
