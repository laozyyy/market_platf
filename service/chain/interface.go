package chain

import (
	"big_market/common/constant"
	"big_market/model"
)

type ActivityChain interface {
	AppendNext(next *ActivityChain) *ActivityChain
	Action(activity model.RaffleActivity, sku model.RaffleActivitySku, count model.RaffleActivityCount) error
}

type LogicChain interface {
	AppendNext(next *LogicChain) *LogicChain
	Logic(userID string, strategyID int64) (int, string, error)
}

var LogicChainGroup = map[string]LogicChain{
	constant.RuleBlacklist: &BlacklistChain{},
	constant.RuleWeight:    &WeightChain{},
	constant.RuleDefault:   &DefaultChain{},
}

var ActivityChainGroup = map[string]ActivityChain{
	constant.ActivityBase:  &ActivityBaseChain{},
	constant.ActivityStock: &ActivityStockChain{},
}
