package logic_chain

import (
	"big_market/common"
	"big_market/database"
	"big_market/reposity"
	"sort"
	"strconv"
	"strings"
)

type WeightChain struct {
	nextChain LogicChain
}

func (w *WeightChain) Next() *LogicChain {
	return &w.nextChain
}

func (w *WeightChain) AppendNext(next *LogicChain) *LogicChain {
	w.nextChain = *next
	return next
}

func (w *WeightChain) Logic(userID string, strategyID int64) (int, error) {
	common.Log.Infof("责任链：权重过滤, userId:%v strategyId:%v", userID, strategyID)
	ruleValue, err := database.QueryStrategyRuleValue(nil, strconv.FormatInt(strategyID, 10), common.RuleWeight, 0)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, nil
	}
	scoreAwardIDMap, keys, err := getScoreAwardIDMapAndSortedKeys(ruleValue)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, nil
	}
	//userScore := database.QueryUserScore(nil, rule.UserID)
	var userScore int64 = 2000
	for _, key := range keys {
		if userScore < key {
			continue
		}
		common.Log.Infof("当前积分级别: %d", key)
		weightValue := scoreAwardIDMap[key]
		awardID := reposity.GetRandomAwardIdByWeight(strconv.FormatInt(strategyID, 10), weightValue)
		return awardID, nil
	}
	next := *(w.Next())
	awardID, err := next.Logic(userID, strategyID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return 0, nil
	}
	return awardID, nil
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
