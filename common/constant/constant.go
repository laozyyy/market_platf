package constant

const (
	StrategyAwardKey     = "big_market_strategy_award_key_"
	StrategyRateTableKey = "big_market_strategy_rate_table_key_"
	StrategyRateRangeKey = "big_market_strategy_rate_range_key_"
	StrategyKey          = "big_market_strategy_key_"
	RuleTreeVOKey        = "rule_tree_vo_key_"

	Split = ","
	Space = " "
	COLON = ":"

	RuleWeight    = "rule_weight"
	RuleBlacklist = "rule_blacklist"
	RuleLuckAward = "rule_luck_award"
	RuleLock      = "rule_lock"
	RuleDefault   = "default"

	LockTreeNode      = "rule_lock"
	LuckAwardTreeNode = "rule_luck_award"
	StockNode         = "rule_stock"

	// Allow 跳过规则
	Allow = "ALLOW"
	// TakeOver 不跳过
	TakeOver = "TAKE_OVER"
)

var BeforeRules = []string{
	RuleWeight,
	RuleBlacklist,
}
var CenterRules = []string{
	RuleLuckAward,
	RuleLock,
}

// RuleWeight 不同规则优先级
var WeightOfRules = map[string]int{
	RuleBlacklist: 1,
	RuleWeight:    2,
}
