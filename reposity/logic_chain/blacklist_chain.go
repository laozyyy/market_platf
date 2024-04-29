package logic_chain

import (
	"big_market/common"
	"big_market/database"
	"strconv"
	"strings"
)

type BlacklistChain struct {
	nextChain LogicChain
}

func (b *BlacklistChain) Next() *LogicChain {
	return &b.nextChain
}

func (b *BlacklistChain) AppendNext(next *LogicChain) *LogicChain {
	b.nextChain = *next
	return next
}

func (b *BlacklistChain) Logic(userID string, strategyID int64) (int, error) {
	common.Log.Infof("责任链：黑名单过滤, userId:%v strategyId:%v", userID, strategyID)
	ruleValue, err := database.QueryStrategyRuleValue(nil, strconv.FormatInt(strategyID, 10), common.RuleBlacklist, 0)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, err
	}
	split := strings.Split(ruleValue, common.COLON)
	blackAwardID := split[0]
	blackAwardIDInt, err := strconv.ParseInt(blackAwardID, 10, 64)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, err
	}
	userIDs := strings.Split(split[1], common.Split)
	for _, blackUserID := range userIDs {
		if userID == blackUserID {
			return int(blackAwardIDInt), nil
		}
	}
	// 不在黑名单,责任链继续
	next := *(b.Next())
	awardID, err := next.Logic(userID, strategyID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, err
	}
	return awardID, nil
}
