package chain

import (
	"big_market/common/constant"
)

type LogicChain interface {
	AppendNext(next *LogicChain) *LogicChain
	Logic(userID string, strategyID int64) (int, string, error)
}

var ChainGroup = map[string]LogicChain{
	constant.RuleBlacklist: &BlacklistChain{},
	constant.RuleWeight:    &WeightChain{},
	constant.RuleDefault:   &DefaultChain{},
}
