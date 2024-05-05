package reposity

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"slices"
	"strconv"
	"strings"
)

// AssembleLotteryStrategyWithRules 装配策略，生成奖品表
func AssembleLotteryStrategyWithRules(strategyID int64) error {
	// 由于strategy_award中的surplus_count一直在更新，使用缓存会使后续缓存surplus_count时不一致
	//strategyAwardList, err := getStrategyAwardList(strategyID)
	strategyAwardList, err := database.QueryStrategyAwardListByStrategyId(nil, strategyID)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	// 缓存
	for _, award := range strategyAwardList {
		err = cacheStrategyAwardCount(strategyID, award.AwardId, award.AwardCountSurplus)
		if err != nil {
			log.Errorf("error: %v", err)
			return err
		}
	}
	err = AssembleLotteryStrategy(strconv.FormatInt(strategyID, 10), strategyAwardList)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	strategy, err := GetStrategyByStrategyID(strategyID)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	models := strings.Split(strategy.RuleModels, constant.Split)
	// todo 不只有权重
	if !slices.Contains(models, "rule_weight") {
		return nil
	}
	result, err := database.QueryStrategyRulesByRuleModel(nil, strategyID, constant.RuleWeight)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	weightValues, err := result.GetWeightValues()
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	for key, awardIDs := range weightValues {
		strategyAwardListClone := make([]*model.StrategyAward, 0)
		for _, award := range strategyAwardList {
			if slices.Contains(awardIDs, award.AwardId) {
				strategyAwardListClone = append(strategyAwardListClone, award)
			}
		}
		err = AssembleLotteryStrategy(fmt.Sprintf("%d_%s", strategyID, key), strategyAwardListClone)
		if err != nil {
			log.Errorf("error: %v", err)
			return err
		}
	}
	return nil
}

func AssembleLotteryStrategy(strategyID string, strategyAwardList []*model.StrategyAward) error {
	var (
		totalRate         float64
		minRate           = math.MaxFloat64
		AwardSearchTables []int
	)

	for _, strategyAward := range strategyAwardList {
		totalRate += strategyAward.AwardRate
		minRate = math.Min(minRate, strategyAward.AwardRate)
	}
	rateRange := math.Ceil(totalRate / minRate)
	for _, strategyAward := range strategyAwardList {
		for i := 0; i < int(math.Ceil(strategyAward.AwardRate*rateRange)); i++ {
			AwardSearchTables = append(AwardSearchTables, strategyAward.AwardId)
		}
	}
	rand.Shuffle(len(AwardSearchTables), func(i int, j int) {
		AwardSearchTables[i], AwardSearchTables[j] = AwardSearchTables[j], AwardSearchTables[i]
	})
	ShuffleAwardSearchTables := make(map[int]int)
	for i, awardID := range AwardSearchTables {
		ShuffleAwardSearchTables[i] = awardID
	}
	err := cacheAwardSearchTables(strategyID, len(ShuffleAwardSearchTables), ShuffleAwardSearchTables)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	return nil
}

// GetRandomAwardIdByWeight 根据权重随机获取奖品
func GetRandomAwardIdByWeight(strategyID string, weight string) (int, error) {
	log.Infof("策略: %s", strategyID)
	log.Infof("权重key: %s", weight)
	rateRange, err := getRateRange(fmt.Sprintf("%s_%s", strategyID, weight))
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	if rateRange <= 0 {
		log.Errorf("RateRange error: %d", rateRange)
		return 0, errors.New("RateRange error")
	}
	random := rand.Intn(rateRange)
	log.Infof("random: %d", random)
	awardID, err := getAwardID(strategyID, random)
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	return awardID, nil
}

// GetRandomAwardId 随机获取奖品
func GetRandomAwardId(strategyID string) (int, error) {
	log.Infof("策略: %s", strategyID)
	rateRange, err := getRateRange(strategyID)
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	if rateRange <= 0 {
		log.Errorf("RateRange error: %d", rateRange)
		return 0, errors.New("RateRange error")
	}
	random := rand.Intn(rateRange)
	log.Infof("random: %d", random)
	awardID, err := getAwardID(strategyID, random)
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	return awardID, nil
}
