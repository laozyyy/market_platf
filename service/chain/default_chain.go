package chain

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/service/reposity"
	"strconv"
)

type DefaultChain struct {
	nextChain LogicChain
}

func (d *DefaultChain) AppendNext(next *LogicChain) *LogicChain {
	d.nextChain = *next
	return next
}

func (d *DefaultChain) Logic(userID string, strategyID int64) (int, string, error) {
	log.Infof("责任链：默认处理")
	awardId, err := reposity.GetRandomAwardId(strconv.FormatInt(strategyID, 10))
	if err != nil {
		return 0, "", err
	}
	return awardId, constant.RuleDefault, nil
}
