package service

import (
	"big_market/common"
	"big_market/database"
	"big_market/model"
	"big_market/model/logic_filter"
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

	strategy, err := GetStrategyByStrategyID(factor.StrategyID)
	common.Log.Infof("策略: %+v", strategy)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	// 抽奖前检查规则
	beforeEntity, err := doCheckRaffleBeforeLogic(factor, strategy.GetRuleModels())
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	// 奖品表的装配与抽奖的执行分开
	if beforeEntity.Code != common.Allow {
		if beforeEntity.RuleModel == common.RuleBlacklist {
			// 黑名单直接返回固定商品
			common.Log.Infof("命中黑名单 userID: %v", factor.UserID)
			return &model.RaffleAward{
				StrategyID: factor.StrategyID,
				AwardID:    beforeEntity.AwardID,
			}, nil
		} else if beforeEntity.RuleModel == common.RuleWeight {
			weightValueKey := beforeEntity.WeightValueKey
			common.Log.Infof("命中权重，权重key: %s", weightValueKey)
			awardId := GetRandomAwardIdByWeight(strconv.FormatInt(factor.StrategyID, 10), weightValueKey)
			return &model.RaffleAward{
				StrategyID: factor.StrategyID,
				AwardID:    awardId,
			}, nil
		}
	}
	common.Log.Infof("未命中抽奖前规则")
	awardId := GetRandomAwardId(strconv.FormatInt(factor.StrategyID, 10))
	ruleModelStr, err := database.QueryStrategyAwardRuleModel(nil, factor.StrategyID, int64(awardId))
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	ruleModels := strings.Split(ruleModelStr, common.Split)
	centerEntity, err := doCheckRaffleCenterLogic(model.RaffleFactor{
		UserID:     factor.UserID,
		StrategyID: factor.StrategyID,
		AwardID:    awardId,
	}, ruleModels)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	if centerEntity.Code != common.Allow {
		if centerEntity.RuleModel == common.RuleLuckAward {
			common.Log.Infof("命中幸运奖，awardID: %d", centerEntity.AwardID)
			return &model.RaffleAward{
				StrategyID: factor.StrategyID,
				AwardID:    centerEntity.AwardID,
			}, nil
		}
	}

	return &model.RaffleAward{
		StrategyID: factor.StrategyID,
		AwardID:    awardId,
	}, nil

}

func doCheckRaffleBeforeLogic(factor model.RaffleFactor, rules []string) (*model.RaffleRuleActionEntity, error) {
	allowedBeforeEntity := &model.RaffleRuleActionEntity{
		Info: common.Allow,
		Code: common.Allow,
	}
	if rules == nil {
		return allowedBeforeEntity, nil
	}
	// 优先过滤黑名单
	if slices.Contains(rules, common.RuleBlacklist) {
		filterRule := logic_filter.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  common.RuleBlacklist,
		}
		beforeEntity, err := logic_filter.LogicFilterGroup[common.RuleBlacklist](filterRule)
		if err != nil {
			common.Log.Errorf("err: %v", err)
			return nil, err
		}
		if beforeEntity.Code != common.Allow {
			return beforeEntity, nil
		}
	}
	// 剩下依次处理
	for _, rule := range rules {
		if rule == common.RuleBlacklist {
			continue
		}
		filterRule := logic_filter.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  rule,
		}
		beforeEntity, err := logic_filter.LogicFilterGroup[rule](filterRule)
		if err != nil {
			common.Log.Errorf("err: %v", err)
			return nil, err
		}
		if beforeEntity.Code != common.Allow {
			return beforeEntity, nil
		}
	}
	// 未命中任何规则
	return allowedBeforeEntity, nil
}

func doCheckRaffleCenterLogic(factor model.RaffleFactor, rules []string) (*model.RaffleRuleActionEntity, error) {
	allowedBeforeEntity := &model.RaffleRuleActionEntity{
		Info: common.Allow,
		Code: common.Allow,
	}
	if rules == nil {
		return allowedBeforeEntity, nil
	}

	for _, rule := range rules {
		if slices.Contains(common.BeforeRules, rule) {
			continue
		}
		filterRule := logic_filter.FilterRule{
			UserID:     factor.UserID,
			StrategyID: strconv.FormatInt(factor.StrategyID, 10),
			RuleModel:  rule,
			AwardID:    factor.AwardID,
		}
		centerEntity, err := logic_filter.LogicFilterGroup[rule](filterRule)
		if err != nil {
			common.Log.Errorf("err: %v", err)
			return nil, err
		}
		if centerEntity.Code != common.Allow {
			return centerEntity, nil
		}
	}
	return allowedBeforeEntity, nil
}
