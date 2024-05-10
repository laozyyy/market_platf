package service

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/model"
	"big_market/service/chain"
	"big_market/service/reposity"
	"errors"
	"strings"
)

func PerformRaffle(factor model.RaffleFactor) (*model.RaffleAward, error) {
	// 如果构造RaffleFactor时没有给string类型参数赋值，会自动赋值为零值 "" ，指针类型会赋值nil

	// 1.参数校验
	if factor.UserID == "" || factor.StrategyID == 0 {
		return nil, errors.New("parameter error")
	}

	// 2.责任链-抽奖前处理
	awardID, ruleModel, err := raffleLogicChain(factor.UserID, factor.StrategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	// 触发规则直接返回
	if ruleModel != constant.RuleDefault {
		return &model.RaffleAward{
			StrategyID: factor.StrategyID,
			AwardID:    awardID,
		}, nil
	}

	// 3.规则树-抽奖中、后处理
	awardID, ruleValue, err := raffleLogicTree(factor.UserID, factor.StrategyID, awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}

	return &model.RaffleAward{
		StrategyID:  factor.StrategyID,
		AwardConfig: ruleValue,
		AwardID:     awardID,
	}, nil

}

func raffleLogicChain(userID string, strategyID int64) (int, string, error) {
	logicChain, err := openLogicChain(strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	awardID, ruleModel, err := logicChain.Logic(userID, strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	return awardID, ruleModel, nil
}

func raffleLogicTree(userID string, strategyID int64, awardID int) (int, string, error) {
	ruleModelStr, err := reposity.GetStrategyAwardRuleValue(strategyID, awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	ruleModels := strings.Split(ruleModelStr, constant.Split)
	treeID := ""
	for _, ruleModel := range ruleModels {
		if strings.HasPrefix(ruleModel, "tree") {
			treeID = ruleModel
		}
	}
	if treeID == "" {
		log.Infof("此策略无规则树，strategyID: %v", strategyID)
		return awardID, "", nil
	}
	ruleTree, err := reposity.GetTreeByTreeID(treeID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	engine := TreeEngine{*ruleTree}
	awardVO, err := engine.Process(userID, strategyID, awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", err
	}
	return awardVO.AwardId, awardVO.RuleValue, nil
}

func openLogicChain(strategyID int64) (chain.LogicChain, error) {
	strategy, err := reposity.GetStrategyByStrategyID(strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	ruleModels := strategy.GetRuleModels()
	if ruleModels == nil {
		log.Infof("无责任链, strategyID: %v", strategyID)
		return chain.LogicChainGroup[constant.RuleDefault], nil
	}
	log.Infof("ruleModels: %+v", ruleModels)
	chainHead := chain.LogicChainGroup[ruleModels[0]]
	current := chainHead
	for i := 1; i < len(ruleModels); i++ {
		chain := chain.LogicChainGroup[ruleModels[i]]
		current = *(current.AppendNext(&chain))
	}
	chain := chain.LogicChainGroup[constant.RuleDefault]
	current.AppendNext(&chain)
	return chainHead, nil
}
