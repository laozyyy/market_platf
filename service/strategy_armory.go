package service

import (
	"big_market/common"
	"math"
	"math/rand"
)

func AssembleLotteryStrategy(strategyID int64) error {
	var (
		totalRate         float64
		minRate           = math.MaxFloat64
		AwardSearchTables []int
	)
	strategyAwardList, err := getStrategyAwardList(strategyID)
	if err != nil {
		common.Log.Errorf("error: %v", err)
		return err
	}

	for _, strategyAward := range strategyAwardList {
		totalRate += strategyAward.AwardRate
		minRate = math.Min(minRate, strategyAward.AwardRate)
	}
	rateRange := math.Ceil(totalRate / minRate)
	// 乱序后续加
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
	err = saveAwardSearchTables(strategyID, len(ShuffleAwardSearchTables), ShuffleAwardSearchTables)
	return nil
}

func GetRandomAwardId(strategyID int64) int {
	rateRange := getRateRange(strategyID)
	random := rand.Intn(rateRange)
	common.Log.Infof("random: %d", random)
	return getAwardID(strategyID, random)
}
