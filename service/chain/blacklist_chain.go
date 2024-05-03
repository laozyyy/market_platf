package chain

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"strconv"
	"strings"
)

type BlacklistChain struct {
	nextChain LogicChain
}

func (b *BlacklistChain) AppendNext(next *LogicChain) *LogicChain {
	b.nextChain = *next
	return next
}

func (b *BlacklistChain) Logic(userID string, strategyID int64) (int, string, error) {
	log.Infof("责任链：黑名单过滤, userId:%v strategyId:%v", userID, strategyID)
	ruleValue, err := database.QueryStrategyRuleValue(nil, strconv.FormatInt(strategyID, 10), constant.RuleBlacklist, 0)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	if ruleValue != "" {
		split := strings.Split(ruleValue, constant.COLON)
		blackAwardID := split[0]
		blackAwardIDInt, err := strconv.ParseInt(blackAwardID, 10, 64)
		if err != nil {
			log.Errorf("err: %v", err)
			return 0, "", err
		}
		userIDs := strings.Split(split[1], constant.Split)
		for _, blackUserID := range userIDs {
			if userID == blackUserID {
				log.Infof("触发黑名单, userId:%v", userID)
				return int(blackAwardIDInt), constant.RuleBlacklist, nil
			}
		}
		log.Infof("不在黑名单, userId:%v", userID)
	} else {
		log.Infof("无黑名单, strategyID:%v", strategyID)
	}
	// 不在黑名单或无黑名单,责任链继续
	awardID, ruleModel, err := b.nextChain.Logic(userID, strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	return awardID, ruleModel, nil
}
