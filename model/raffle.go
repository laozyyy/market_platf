package model

type RaffleFactor struct {
	UserID     string `json:"user_id"`
	StrategyID int64  `json:"strategy_id"`
	AwardID    int    `json:"award_id"`
}

type RaffleAward struct {
	StrategyID  int64  `json:"strategy_id"`  // 策略ID
	AwardID     int    `json:"award_id"`     // 奖品ID
	AwardKey    string `json:"award_key"`    // 奖品对接标识 - 每一个都是一个对应的发奖策略
	AwardConfig string `json:"award_config"` // 奖品配置信息
	AwardDesc   string `json:"award_desc"`   // 奖品内容描述
}

type RaffleRuleActionEntity struct {
	Info           string `json:"info"`
	Code           string `json:"code"`
	RuleModel      string `json:"rule_model"`
	StrategyID     string `json:"strategy_id"`
	WeightValueKey string `json:"weight_value_key"`
	AwardID        int    `json:"award_id"`
}
