package model

import (
	"time"
)

const (
	STRATEGY_AWARD_KEY      = "big_market_strategy_award_key_"
	STRATEGY_RATE_TABLE_KEY = "big_market_strategy_rate_table_key_"
	STRATEGY_RATE_RANGE_KEY = "big_market_strategy_rate_range_key_"
)

type Award struct {
	ID          int       `json:"id"`
	AwardID     int       `json:"awardId"`
	AwardKey    string    `json:"awardKey"`
	AwardConfig string    `json:"awardConfig"`
	AwardDesc   string    `json:"awardDesc"`
	CreateTime  time.Time `json:"createTime"`
	UpdateTime  time.Time `json:"updateTime"`
}

type Strategy struct {
	ID           int       `json:"id"`
	StrategyID   int       `json:"strategyId"`
	StrategyDesc string    `json:"strategyDesc"`
	CreateTime   time.Time `json:"createTime"`
	UpdateTime   time.Time `json:"updateTime"`
}

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

type StrategyRule struct {
	ID         int       `json:"id"`         // 自增ID
	StrategyID int       `json:"strategyId"` // 抽奖策略ID
	AwardID    int       `json:"awardId"`    // 抽奖奖品ID【规则类型为策略，则不需要奖品ID】
	RuleType   int       `json:"ruleType"`   // 抽象规则类型；1-策略规则、2-奖品规则
	RuleModel  string    `json:"ruleModel"`  // 抽奖规则类型【rule_random - 随机值计算、rule_lock - 抽奖几次后解锁、rule_luck_award - 幸运奖(兜底奖品)】
	RuleValue  string    `json:"ruleValue"`  // 抽奖规则比值
	RuleDesc   string    `json:"ruleDesc"`   // 抽奖规则描述
	CreateTime time.Time `json:"createTime"` // 创建时间
	UpdateTime time.Time `json:"updateTime"` // 更新时间
}
