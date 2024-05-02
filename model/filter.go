package model

// FilterRule 待过滤的规则
type FilterRule struct {
	UserID     string `json:"user_id"`
	StrategyID string `json:"strategy_id"`
	AwardID    int    `json:"award_id"`
	RuleModel  string `json:"rule_model"`
}
