package service

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/model"
	"big_market/mq"
	"big_market/service/reposity"
	"errors"
	"strconv"
	"strings"
)

var LogicTreeNodeGroup = map[string]TreeNode{
	constant.LockTreeNode:      LockTreeNode{},
	constant.LuckAwardTreeNode: LuckAwardTreeNode{},
	constant.StockNode:         StockTreeNode{},
}

type TreeNode interface {
	Logic(userID string, strategyID int64, awardID int, ruleValue string) (code string, strategyAwardVO *model.StrategyAwardVO, err error)
}

type LockTreeNode struct {
}

type LuckAwardTreeNode struct {
}

type StockTreeNode struct {
}

func (l LockTreeNode) Logic(userID string, strategyID int64, awardID int, ruleValueStr string) (code string, strategyAwardVO *model.StrategyAwardVO, err error) {
	log.Infof("决策树-积分解锁")
	userRaffleCount := 10
	ruleValue, err := strconv.ParseInt(ruleValueStr, 10, 64)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", nil, err
	}
	if int64(userRaffleCount) > ruleValue {
		log.Infof("决策树-积分解锁-放行: 次数为%d, 阈值为:%d", userRaffleCount, ruleValue)
		return constant.Allow, &model.StrategyAwardVO{
			AwardId:   awardID,
			RuleValue: ruleValueStr,
		}, nil
	}
	log.Infof("决策树-积分解锁-接管: 次数为%d, 阈值为:%d", userRaffleCount, ruleValue)
	return constant.TakeOver, &model.StrategyAwardVO{
		AwardId:   awardID,
		RuleValue: ruleValueStr,
	}, nil
}

func (s StockTreeNode) Logic(userID string, strategyID int64, awardID int, ruleValue string) (code string, strategyAwardVO *model.StrategyAwardVO, err error) {
	log.Infof("决策树-库存")
	success, err := reposity.DescStrategyAwardCountCache(strategyID, awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", nil, err
	}
	if success {
		// 写入延迟队列，延迟消费更新数据库记录
		err = mq.SendUpdateAwardCountMessage(strategyID, awardID)
		if err != nil {
			log.Errorf("err: %v", err)
			return "", nil, err
		}
		return constant.TakeOver, &model.StrategyAwardVO{
			AwardId:   awardID,
			RuleValue: ruleValue,
		}, nil
	}
	return constant.Allow, strategyAwardVO, nil
}

func (l LuckAwardTreeNode) Logic(userID string, strategyID int64, awardID int, ruleValue string) (code string, strategyAwardVO *model.StrategyAwardVO, err error) {
	log.Infof("决策树-幸运奖")
	// 数据: 101:1,100
	split := strings.Split(ruleValue, constant.Colon)
	if len(split) == 0 {
		return constant.Allow, nil, errors.New("luck award not configured")
	}
	luckAwardID := split[0]
	luckAwardIDInt, err := strconv.ParseInt(luckAwardID, 10, 64)
	log.Infof("决策树-幸运奖-awardID: %v", awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", nil, err
	}
	strategyAwardVO = &model.StrategyAwardVO{
		AwardId: int(luckAwardIDInt),
	}
	if len(split) > 1 {
		strategyAwardVO.RuleValue = split[1]
	} else {
		log.Infof("决策树-幸运奖-无积分，strategyID: %v", strategyID)
	}
	return constant.TakeOver, strategyAwardVO, err
}
