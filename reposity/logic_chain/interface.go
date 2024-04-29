package logic_chain

import "big_market/common"

type LogicChain interface {
	Next() *LogicChain
	AppendNext(next *LogicChain) *LogicChain
	Logic(userID string, strategyID int64) (int, error)
}

var ChainGroup = map[string]LogicChain{
	common.RuleBlacklist: &BlacklistChain{},
	common.RuleWeight:    &WeightChain{},
	common.RuleDefault:   &DefaultChain{},
}
