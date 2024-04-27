package service

import (
	"big_market/cache"
	"big_market/common"
	"big_market/database"
	"big_market/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

// 先从缓存获取
func getStrategyAwardList(strategyID int64) ([]*model.StrategyAward, error) {
	ctx := context.Background()
	key := model.StrategyAwardKey + strconv.FormatInt(strategyID, 10)
	result, err := cache.Client.Get(ctx, key).Result()
	if !errors.Is(redis.Nil, err) && result != "" {
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
	err = cache.Client.Set(ctx, key, strategyAwardsStr, 0).Err()
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	return strategyAwards, nil
}

func saveAwardSearchTables(strategyID string, rateRange int, rateRangeTable map[int]int) error {
	ctx := context.Background()
	rateRangeKey := fmt.Sprintf("%s%s", model.StrategyRateRangeKey, strategyID)
	rateTableKey := fmt.Sprintf("%s%s", model.StrategyRateTableKey, strategyID)

	cache.Client.Set(ctx, rateRangeKey, rateRange, 0)

	oldTable, err := cache.Client.HGetAll(ctx, rateTableKey).Result()
	if err != nil {
		common.Log.Infof("err: %v", err)
	}
	for k, v := range rateRangeTable {
		oldTable[strconv.Itoa(k)] = strconv.Itoa(v)
	}
	err = cache.Client.HSet(ctx, rateTableKey, oldTable).Err()
	if err != nil {
		common.Log.Infof("err: %v", err)
		return err
	}
	return nil
}

func getRateRange(strategyID int64) int {
	ctx := context.Background()
	rateRange, err := cache.Client.Get(ctx, fmt.Sprintf("%s%d", model.StrategyRateRangeKey, strategyID)).Result()
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
	awardID, err := cache.Client.HGet(ctx, fmt.Sprintf("%s%d", model.StrategyRateTableKey, strategyID), strconv.Itoa(rateKey)).Result()
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	awardIDInt, err := strconv.Atoi(awardID)
	if err != nil {
		common.Log.Errorf("error: %v", err)
	}
	return awardIDInt
}
