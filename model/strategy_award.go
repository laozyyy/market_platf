package model

import "time"

type StrategyAward struct {
	ID                int       `json:"id"`                // 自增ID
	StrategyID        int       `json:"strategyId"`        // 抽奖策略ID
	AwardId           int       `json:"awardId"`           // 抽奖奖品ID - 内部流转使用
	AwardTitle        string    `json:"awardTitle"`        // 抽奖奖品标题
	AwardSubtitle     string    `json:"awardSubtitle"`     // 抽奖奖品副标题
	AwardCount        int       `json:"awardCount"`        // 奖品库存总量
	AwardCountSurplus int       `json:"awardCountSurplus"` // 奖品库存剩余
	AwardRate         float64   `json:"awardRate"`         // 奖品中奖概率
	RuleModels        string    `json:"ruleModels"`        // 规则模型，rule配置的模型同步到此表，便于使用
	Sort              int       `json:"sort"`              // 排序
	CreateTime        time.Time `json:"createTime"`        // 创建时间
	UpdateTime        time.Time `json:"updateTime"`        // 修改时间
}
