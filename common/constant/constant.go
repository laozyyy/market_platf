package constant

const (
	StrategyAwardKey          = "big_market_strategy_award_key_"
	StrategyRateTableKey      = "big_market_strategy_rate_table_key_"
	StrategyRateRangeKey      = "big_market_strategy_rate_range_key_"
	StrategyKey               = "big_market_strategy_key_"
	RuleTreeVOKey             = "rule_tree_vo_key_"
	StrategyAwardCountKey     = "strategy_award_count_key_"
	StrategyAwardRuleValueKey = "strategy_award_rule_value_key_"
	ActivitySkuStockCountKey  = "activity_sku_stock_count_key"

	ActivityKey      = "big_market_activity_key_"
	ActivityCountKey = "big_market_activity_count_key_"

	Split = ","
	Space = " "
	Colon = ":"

	RuleWeight    = "rule_weight"
	RuleBlacklist = "rule_blacklist"
	RuleDefault   = "default"

	ActivityBase  = "activity_base"
	ActivityStock = "activity_stock"

	ActivityCreate = "create"
	ActivityOpen   = "open"
	ActivityClose  = "close"

	LockTreeNode      = "rule_lock"
	LuckAwardTreeNode = "rule_luck_award"
	StockNode         = "rule_stock"

	// Allow 跳过规则
	Allow = "ALLOW"
	// TakeOver 不跳过
	TakeOver = "TAKE_OVER"

	// 实际上只有交换机延迟
	DelayQueueName = "delay_queue"

	DelayExchangeName  = "delayed-exchange"
	NormalExchangeName = "normal-exchange"

	UpdateStrategyAwardCountTopic = "update_strategy_award_count_topic"
	SkuCountZeroTopic             = "sku_count_zero_topic"
	UpdateSkuCountTopic           = "update_sku_count_topic"
	UpdateSkuCountQueue           = "update_sku_count_queue"

	OrderCompleted = "order_completed"
)

// RuleWeight 不同规则优先级
var WeightOfRules = map[string]int{
	RuleBlacklist: 1,
	RuleWeight:    2,
}
