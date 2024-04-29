package logic_filter

import (
	"big_market/common"
	"big_market/database"
	"big_market/model"
	"errors"
	"sort"
	"strconv"
	"strings"
)

var LogicFilterGroup = map[string]func(rule FilterRule) (*model.RaffleRuleActionEntity, error){
	common.RuleBlacklist: BlackListFilter,
	common.RuleWeight:    WeightFilter,
	common.RuleLuckAward: LuckAwardFilter,
	common.RuleLock:      LockFilter,
}

// FilterRule 待过滤的规则
type FilterRule struct {
	UserID     string `json:"user_id"`
	StrategyID string `json:"strategy_id"`
	AwardID    int    `json:"award_id"`
	RuleModel  string `json:"rule_model"`
}

func BlackListFilter(rule FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	common.Log.Infof("黑名单过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	split := strings.Split(ruleValue, common.COLON)
	awardID := split[0]
	awardIDInt, err := strconv.Atoi(awardID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	userIDs := strings.Split(split[1], common.Split)
	for _, blackUserID := range userIDs {
		if rule.UserID == blackUserID {
			return &model.RaffleRuleActionEntity{
				Info:           common.TakeOver,
				Code:           common.TakeOver,
				RuleModel:      common.RuleBlacklist,
				StrategyID:     rule.StrategyID,
				WeightValueKey: "",
				AwardID:        awardIDInt,
			}, nil
		}
	}
	// 不在黑名单
	return &model.RaffleRuleActionEntity{
		Info: common.Allow,
		Code: common.Allow,
	}, nil

}

func WeightFilter(rule FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	common.Log.Infof("权重过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	scoreAwardIDMap, keys, err := getScoreAwardIDMapAndSortedKeys(ruleValue)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	//userScore := database.QueryUserScore(nil, rule.UserID)
	var userScore int64 = 5000
	for _, key := range keys {
		if userScore < key {
			continue
		}
		common.Log.Infof("当前积分级别: %d", key)
		ruleValue := scoreAwardIDMap[key]
		return &model.RaffleRuleActionEntity{
			Info:           common.TakeOver,
			Code:           common.TakeOver,
			RuleModel:      common.RuleWeight,
			StrategyID:     rule.StrategyID,
			WeightValueKey: ruleValue,
		}, nil
	}
	return &model.RaffleRuleActionEntity{
		Info: common.Allow,
		Code: common.Allow,
	}, nil
}

func LuckAwardFilter(rule FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	common.Log.Infof("幸运奖过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	_, err = database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	if err != nil && !errors.Is(err, common.NoDataErr) {
		common.Log.Errorf("err: %v", err)
		return
	}
	return &model.RaffleRuleActionEntity{
		Info: common.Allow,
		Code: common.Allow,
	}, nil
}

func LockFilter(rule FilterRule) (result *model.RaffleRuleActionEntity, err error) {
	common.Log.Infof("抽奖次数解锁过滤, userId:%v strategyId:%v ruleModel:%v", rule.UserID, rule.StrategyID, rule.RuleModel)
	_, err = database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, rule.AwardID)
	userCount := 1

	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	ruleValue, err := database.QueryStrategyRuleValue(nil, rule.StrategyID, rule.RuleModel, 0)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	ruleValueInt, err := strconv.ParseInt(ruleValue, 10, 64)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return
	}
	if int64(userCount) > ruleValueInt {
		common.Log.Infof("命中抽奖次数解锁规则")
		return &model.RaffleRuleActionEntity{
			Info: common.Allow,
			Code: common.Allow,
		}, nil
	}
	return &model.RaffleRuleActionEntity{
		Info: common.TakeOver,
		Code: common.TakeOver,
	}, nil

}

func getScoreAwardIDMapAndSortedKeys(ruleValue string) (map[int64]string, []int64, error) {
	ruleValueGroups := strings.Split(ruleValue, common.Space)
	scoreAwardIDMap := make(map[int64]string)
	for _, ruleValueGroup := range ruleValueGroups {
		split := strings.Split(ruleValueGroup, common.COLON)
		score, err := strconv.ParseInt(split[0], 10, 64)
		if err != nil {
			common.Log.Errorf("err: %v", err)
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
