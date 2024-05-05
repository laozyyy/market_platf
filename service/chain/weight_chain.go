package chain

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/service/reposity"
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

func (w *WeightChain) Logic(userID string, strategyID int64) (int, string, error) {
	log.Infof("责任链：权重过滤, userId:%v strategyId:%v", userID, strategyID)
	ruleValue, err := database.QueryStrategyRuleValue(nil, strconv.FormatInt(strategyID, 10), constant.RuleWeight, 0)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", nil
	}
	scoreAwardIDMap, keys, err := getScoreAwardIDMapAndSortedKeys(ruleValue)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", nil
	}
	//userScore := database.QueryUserScore(nil, rule.UserID)
	var userScore int64 = 6000
	for _, key := range keys {
		if userScore < key {
			continue
		}
		log.Infof("当前积分级别: %d", key)
		weightValue := scoreAwardIDMap[key]
		awardID, err := reposity.GetRandomAwardIdByWeight(strconv.FormatInt(strategyID, 10), weightValue)
		if err != nil {
			log.Errorf("err: %v", err)
			return 0, "", nil
		}
		return awardID, constant.RuleWeight, nil
	}
	log.Infof("未触发权重过滤，score: %d", userScore)
	awardID, ruleModel, err := w.nextChain.Logic(userID, strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, "", nil
	}
	return awardID, ruleModel, nil
}

func getScoreAwardIDMapAndSortedKeys(ruleValue string) (map[int64]string, []int64, error) {
	ruleValueGroups := strings.Split(ruleValue, constant.Space)
	scoreAwardIDMap := make(map[int64]string)
	for _, ruleValueGroup := range ruleValueGroups {
		split := strings.Split(ruleValueGroup, constant.Colon)
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
