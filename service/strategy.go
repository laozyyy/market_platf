package service

import (
	"big_market/common"
	"big_market/database"
	"big_market/model"
	"big_market/redis"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
)

// 先从缓存获取
func getStrategyAwardList(strategyID int64) ([]*model.StrategyAward, error) {
	ctx := context.Background()
	key := model.STRATEGY_AWARD_KEY + strconv.FormatInt(strategyID, 10)
	result, err := redis.Client.Get(ctx, key).Result()
	if result != "" {
		strategyAwards := make([]*model.StrategyAward, 0)
		err := json.Unmarshal([]byte(result), &strategyAwards)
		if err != nil {
			common.Log.Errorf("err: %v", err)
			return nil, err
		}
		return strategyAwards, nil
	}
	strategyAwards, err := database.QueryStrategyAwardListByStrategyId(nil, strategyID)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	strategyAwardsStr, err := json.Marshal(strategyAwards)
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	err = redis.Client.Set(ctx, key, strategyAwardsStr, 0).Err()
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	return strategyAwards, nil
}

func saveAwardSearchTables(strategyID int64, rateRange int, rateRangeTable map[int]int) error {
	ctx := context.Background()
	rateRangeKey := fmt.Sprintf("%s%d", model.STRATEGY_RATE_RANGE_KEY, strategyID)
	rateTableKey := fmt.Sprintf("%s%d", model.STRATEGY_RATE_TABLE_KEY, strategyID)

	redis.Client.Set(ctx, rateRangeKey, rateRange, 0)

	oldTable, err := redis.Client.HGetAll(ctx, rateTableKey).Result()
	if err != nil {
		common.Log.Infof("err: %v", err)
	}
	for k, v := range rateRangeTable {
		oldTable[strconv.Itoa(k)] = strconv.Itoa(v)
	}
	err = redis.Client.HSet(ctx, rateTableKey, oldTable).Err()
	if err != nil {
		common.Log.Infof("err: %v", err)
		return err
	}
	return nil
}

func getRateRange(strategyID int64) int {
	ctx := context.Background()
	rateRange, err := redis.Client.Get(ctx, fmt.Sprintf("%s%d", model.STRATEGY_RATE_RANGE_KEY, strategyID)).Result()
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	rateRangeInt, err := strconv.Atoi(rateRange)
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	return rateRangeInt
}

func getAwardID(strategyID int64, rateKey int) int {
	ctx := context.Background()
	awardID, err := redis.Client.HGet(ctx, fmt.Sprintf("%s%d", model.STRATEGY_RATE_TABLE_KEY, strategyID), strconv.Itoa(rateKey)).Result()
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	awardIDInt, err := strconv.Atoi(awardID)
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	return awardIDInt
}
