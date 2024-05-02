package service

import (
	"big_market/common"
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"errors"
	"sort"
	"strconv"
	"strings"
)

var LogicFilterGroup = map[string]func(rule model.FilterRule) (*model.RaffleRuleActionEntity, error){
	constant.RuleBlacklist: BlackListFilter,
	constant.RuleWeight:    WeightFilter,
	constant.RuleLuckAward: LuckAwardFilter,
	constant.RuleLock:      LockFilter,
}

func BlackListFilter(rule model.FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	log.Infof("黑名单过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	split := strings.Split(ruleValue, constant.COLON)
	awardID := split[0]
	awardIDInt, err := strconv.Atoi(awardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	userIDs := strings.Split(split[1], constant.Split)
	for _, blackUserID := range userIDs {
		if rule.UserID == blackUserID {
			return &model.RaffleRuleActionEntity{
				Info:           constant.TakeOver,
				Code:           constant.TakeOver,
				RuleModel:      constant.RuleBlacklist,
				StrategyID:     rule.StrategyID,
				WeightValueKey: "",
				AwardID:        awardIDInt,
			}, nil
		}
	}
	// 不在黑名单
	return &model.RaffleRuleActionEntity{
		Info: constant.Allow,
		Code: constant.Allow,
	}, nil

}

func WeightFilter(rule model.FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	log.Infof("权重过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	scoreAwardIDMap, keys, err := getScoreAwardIDMapAndSortedKeys(ruleValue)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	//userScore := database.QueryUserScore(nil, rule.UserID)
	var userScore int64 = 5000
	for _, key := range keys {
		if userScore < key {
			continue
		}
		log.Infof("当前积分级别: %d", key)
		ruleValue := scoreAwardIDMap[key]
		return &model.RaffleRuleActionEntity{
			Info:           constant.TakeOver,
			Code:           constant.TakeOver,
			RuleModel:      constant.RuleWeight,
			StrategyID:     rule.StrategyID,
			WeightValueKey: ruleValue,
		}, nil
	}
	return &model.RaffleRuleActionEntity{
		Info: constant.Allow,
		Code: constant.Allow,
	}, nil
}

func LuckAwardFilter(rule model.FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	log.Infof("幸运奖过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	_, err = database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil && !errors.Is(err, common.NoDataErr) {
		log.Errorf("err: %v", err)
		return
	}
	return &model.RaffleRuleActionEntity{
		Info: constant.Allow,
		Code: constant.Allow,
	}, nil
}

func LockFilter(rule model.FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	log.Infof("抽奖次数解锁过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	_, err = database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	userCount := 1

	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, 0)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	ruleValueInt, err := strconv.ParseInt(ruleValue, 10, 64)
	if err != nil {
		log.Errorf("err: %v", err)
		return
	}
	if int64(userCount) > ruleValueInt {
		log.Infof("命中抽奖次数解锁规则")
		return &model.RaffleRuleActionEntity{
			Info: constant.Allow,
			Code: constant.Allow,
		}, nil
	}
	return &model.RaffleRuleActionEntity{
		Info: constant.TakeOver,
		Code: constant.TakeOver,
	}, nil

}

func getScoreAwardIDMapAndSortedKeys(ruleValue string) (map[int64]string, []int64, error) {
	ruleValueGroups := strings.Split(ruleValue, constant.Space)
	scoreAwardIDMap := make(map[int64]string)
	for _, ruleValueGroup := range ruleValueGroups {
		split := strings.Split(ruleValueGroup, constant.COLON)
		score, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, nil, err
		}
		scoreAwardIDMap[score] = ruleValueGroup
	}
	keys := make([]int64, 0)
	for key := range scoreAwardIDMap {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] > keys[j]
	})
	return scoreAwardIDMap, keys, nil
}
