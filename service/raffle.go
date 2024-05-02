package service

import (
	"big_market/common"
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"big_market/reposity"
	"big_market/reposity/logic_chain"
	"errors"
	"slices"
	"strconv"
	"strings"
)

func PerformRaffle(factor model.RaffleFactor) (*model.RaffleAward, error) {
	// 如果构造RaffleFactor时没有给string类型参数赋值，会自动赋值为零值 "" ，指针类型会赋值nil
	if factor.UserID == "" || factor.StrategyID == 0 {
		return nil, errors.New("parameter error")
	}

	strategy, err := reposity.GetStrategyByStrategyID(factor.StrategyID)
	log.Infof("策略: %+v", strategy)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	//// 抽奖前检查规则
	//beforeEntity, err := doCheckRaffleBeforeLogic(factor, strategy.GetRuleModels())
	//if err != nil {
	//	common.Log.Errorf("err: %v", err)
	//	return nil, err
	//}
	//// 奖品表的装配与抽奖的执行分开
	//if beforeEntity.Code != common.Allow {
	//	if beforeEntity.RuleModel == common.RuleBlacklist {
	//		// 黑名单直接返回固定商品
	//		common.Log.Infof("命中黑名单 userID: %v", factor.UserID)
	//		return &model.RaffleAward{
	//			StrategyID: factor.StrategyID,
	//			AwardID:    beforeEntity.AwardID,
	//		}, nil
	//	} else if beforeEntity.RuleModel == common.RuleWeight {
	//		weightValueKey := beforeEntity.WeightValueKey
	//		common.Log.Infof("命中权重，权重key: %s", weightValueKey)
	//		awardId := reposity.GetRandomAwardIdByWeight(strconv.FormatInt(factor.StrategyID, 10), weightValueKey)
	//		return &model.RaffleAward{
	//			StrategyID: factor.StrategyID,
	//			AwardID:    awardId,
	//		}, nil
	//	}
	//}
	//common.Log.Infof("未命中抽奖前规则")
	//awardID := reposity.GetRandomAwardId(strconv.FormatInt(factor.StrategyID, 10))

	logicChain, err := openLogicChain(factor.StrategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	awardID, err := logicChain.Logic(factor.UserID, factor.StrategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}

	ruleModelStr, err := database.QueryStrategyAwardRuleModel(nil, factor.StrategyID, int64(awardID))
	if err != nil && !errors.Is(err, common.NoDataErr) {
		log.Errorf("err: %v", err)
		return nil, err
	}
	ruleModels := strings.Split(ruleModelStr, constant.Split)
	centerEntity, err := doCheckRaffleCenterLogic(model.RaffleFactor{
		UserID:     factor.UserID,
		StrategyID: factor.StrategyID,
		AwardID:    awardID,
	}, ruleModels)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	if centerEntity.Code != constant.Allow {
		log.Infof("\"【临时日志】中奖中规则拦截，通过抽奖后规则 rule_luck_award 走兜底奖励。\"")
		if centerEntity.RuleModel == constant.RuleLuckAward {
			log.Infof("命中幸运奖，awardID: %d", centerEntity.AwardID)
			return &model.RaffleAward{
				StrategyID: factor.StrategyID,
				AwardID:    centerEntity.AwardID,
			}, nil
		}
	}

	return &model.RaffleAward{
		StrategyID: factor.StrategyID,
		AwardID:    awardID,
	}, nil

}

func doCheckRaffleBeforeLogic(factor model.RaffleFactor, rules []string) (*model.RaffleRuleActionEntity, error) {
	allowedBeforeEntity := &model.RaffleRuleActionEntity{
		Info: constant.Allow,
		Code: constant.Allow,
	}
	if rules == nil {
		return allowedBeforeEntity, nil
	}
	// 优先过滤黑名单
	if slices.Contains(rules, constant.RuleBlacklist) {
		filterRule := model.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  constant.RuleBlacklist,
		}
		beforeEntity, err := LogicFilterGroup[constant.RuleBlacklist](filterRule)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		if beforeEntity.Code != constant.Allow {
			return beforeEntity, nil
		}
	}
	// 剩下依次处理
	for _, rule := range rules {
		if rule == constant.RuleBlacklist {
			continue
		}
		filterRule := model.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  rule,
		}
		beforeEntity, err := LogicFilterGroup[rule](filterRule)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		if beforeEntity.Code != constant.Allow {
			return beforeEntity, nil
		}
	}
	// 未命中任何规则
	return allowedBeforeEntity, nil
}

func doCheckRaffleCenterLogic(factor model.RaffleFactor, rules []string) (*model.RaffleRuleActionEntity, error) {
	allowedBeforeEntity := &model.RaffleRuleActionEntity{
		Info: constant.Allow,
		Code: constant.Allow,
	}
	if rules == nil || rules[0] == "" {
		return allowedBeforeEntity, nil
	}

	for _, rule := range rules {
		if slices.Contains(constant.BeforeRules, rule) {
			continue
		}
		filterRule := model.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  rule,
			AwardID:    factor.AwardID,
		}
		centerEntity, err := LogicFilterGroup[rule](filterRule)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		if centerEntity.Code != constant.Allow {
			return centerEntity, nil
		}
	}
	return allowedBeforeEntity, nil
}

func openLogicChain(strategyID int64) (logic_chain.LogicChain, error) {
	strategy, err := reposity.GetStrategyByStrategyID(strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	ruleModels := strategy.GetRuleModels()
	log.Infof("ruleModels: %+v", ruleModels)
	if ruleModels == nil {
		return nil, errors.New("no rule models")
	}
	chainHead := logic_chain.ChainGroup[ruleModels[0]]
	current := chainHead
	for i := 1; i < len(ruleModels); i++ {
		chain := logic_chain.ChainGroup[ruleModels[i]]
		current = *(current.AppendNext(&chain))
	}
	chain := logic_chain.ChainGroup[constant.RuleDefault]
	current.AppendNext(&chain)
	return chainHead, nil
}
