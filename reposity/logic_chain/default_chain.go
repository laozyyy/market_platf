package logic_chain

import (
	"big_market/common"
	"big_market/reposity"
	"strconv"
)

type DefaultChain struct {
	nextChain LogicChain
}

func (d *DefaultChain) Next() *LogicChain {
	return &d.nextChain
}

func (d *DefaultChain) AppendNext(next *LogicChain) *LogicChain {
	d.nextChain = *next
	return next
}

func (d DefaultChain) Logic(userID string, strategyID int64) (int, error) {
	common.Log.Infof("责任链：默认处理")
	return reposity.GetRandomAwardId(strconv.FormatInt(strategyID, 10)), nil
}
