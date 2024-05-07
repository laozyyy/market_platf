package reposity

import (
	"big_market/cache"
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

// GetStrategyAwardList 先从缓存获取
func GetStrategyAwardList(strategyID int64) ([]*model.StrategyAward, error) {
	ctx := context.Background()
	key := getStrategyAwardKey(strategyID)
	result, err := cache.Client.Get(ctx, key).Result()
	// todo 改了这里
	if result != "" {
		strategyAwards := make([]*model.StrategyAward, 0)
		err := json.Unmarshal([]byte(result), &strategyAwards)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		return strategyAwards, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("err: %v", err)
		return nil, err
	}
	strategyAwards, err := database.QueryStrategyAwardListByStrategyId(nil, strategyID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	strategyAwardsStr, err := json.Marshal(strategyAwards)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = cache.Client.Set(ctx, key, strategyAwardsStr, 0).Err()
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return strategyAwards, nil
}

func cacheAwardSearchTables(strategyID string, rateRange int, rateRangeTable map[int]int) error {
	ctx := context.Background()
	rateRangeKey := fmt.Sprintf("%s%s", constant.StrategyRateRangeKey, strategyID)
	rateTableKey := fmt.Sprintf("%s%s", constant.StrategyRateTableKey, strategyID)

	err := cache.Client.Set(ctx, rateRangeKey, rateRange, 0).Err()
	if err != nil {
		log.Infof("err: %v", err)
		return err
	}

	oldTable, err := cache.Client.HGetAll(ctx, rateTableKey).Result()
	if err != nil {
		log.Infof("err: %v", err)
	}
	for k, v := range rateRangeTable {
		oldTable[strconv.Itoa(k)] = strconv.Itoa(v)
	}
	err = cache.Client.HSet(ctx, rateTableKey, oldTable).Err()
	if err != nil {
		log.Infof("err: %v", err)
		return err
	}
	return nil
}

func GetStrategyAwardRuleValue(strategyID int64, awardID int) (string, error) {
	key := fmt.Sprintf("%s%d_%d", constant.StrategyAwardRuleValueKey, strategyID, awardID)
	result, err := cache.Client.Get(context.Background(), key).Result()
	if result != "" {
		return result, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Infof("err: %v", err)
		return "", nil
	}
	ruleModelStr, err := database.QueryStrategyAwardRuleModel(nil, strategyID, awardID)
	err = cache.Client.Set(context.Background(), key, ruleModelStr, 0).Err()
	if err != nil {
		log.Infof("err: %v", err)
		return "", nil
	}
	return ruleModelStr, nil
}

func getRateRange(strategyID string) (int, error) {
	ctx := context.Background()
	rateRange, err := cache.Client.Get(ctx, fmt.Sprintf("%s%s", constant.StrategyRateRangeKey, strategyID)).Result()
	if rateRange != "" {
		rateRangeInt, _ := strconv.Atoi(rateRange)
		return rateRangeInt, err
	}
	log.Errorf("error: %v", err)
	return 0, nil
}

func getAwardID(strategyID string, rateKey int) (int, error) {
	ctx := context.Background()
	awardID, err := cache.Client.HGet(ctx, fmt.Sprintf("%s%s", constant.StrategyRateTableKey, strategyID), strconv.Itoa(rateKey)).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	awardIDInt, err := strconv.Atoi(awardID)
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	return awardIDInt, nil
}

func cacheStrategyAwardCount(strategyID int64, awardID int, awardCount int) error {
	key := getStrategyAwardCountKey(strategyID, awardID)
	err := cache.Client.Set(context.Background(), key, awardCount, 0).Err()
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	return nil
}

func DescStrategyAwardCountCache(strategyID int64, awardID int) (bool, error) {
	key := getStrategyAwardCountKey(strategyID, awardID)
	surplus, err := cache.Client.Decr(context.Background(), key).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return false, err
	}
	if surplus < 0 {
		// 恢复库存
		log.Infof("库存为0 awardID: %d awardID: %d", awardID, strategyID)
		err = cache.Client.Set(context.Background(), key, 0, 0).Err()
		if err != nil {
			log.Errorf("error: %v", err)
			return false, err
		}
	}
	if err != nil {
		log.Errorf("error: %v", err)
		return false, err
	}
	// 1. 按照cacheKey decr 后的值，如 99、98、97 和 key 组成为库存锁的key进行使用。
	// 2. 加锁为了兜底，如果后续有恢复库存，手动处理等，也不会超卖。因为所有的可用库存key，都被加锁了。
	lockKey := fmt.Sprintf("%s_%d", key, surplus)
	result, err := cache.Client.SetNX(context.Background(), lockKey, 1, time.Second).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return false, err
	}
	// 这里如果加锁失败，会导致库存减少但是没有抽奖成功，不会超卖
	// 超卖：用户抽奖成功，但实际上无法获得奖品（库存不足）
	// 是否可以优化？恢复库存数
	if !result {
		log.Infof("策略奖品库存加锁失败 %s", lockKey)
		return false, nil
	}
	err = cache.Client.Del(context.Background(), lockKey).Err()
	if err != nil {
		log.Errorf("error: %v", err)
		return false, err
	}
	return true, nil
}

// UpdateStrategyAwardCount 缓存更新数据库
func UpdateStrategyAwardCount(strategyID int64, awardID int) error {
	key := getStrategyAwardCountKey(strategyID, awardID)
	surplus, err := cache.Client.Get(context.Background(), key).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	// 使用database中的连接单例 防止多次执行任务导致连接数过多
	err = database.UpdateStrategyAwardAwardCountSurplus(database.DB, strategyID, awardID, surplus)
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	// 旁路缓存 写库后删缓存
	if err = DeleteCacheString(getStrategyAwardKey(strategyID)); err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	return nil
}

func DeleteCacheString(key string) error {
	result, err := cache.Client.Del(context.Background(), key).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return err
	}
	if result > 0 {
		log.Infof("共删除%v条缓存, key: %s", result, key)
	}
	return nil
}

func getStrategyAwardKey(strategyID int64) string {
	return fmt.Sprintf("%s%d", constant.StrategyAwardKey, strategyID)
}

func getStrategyAwardCountKey(strategyID int64, awardID int) string {
	return fmt.Sprintf("%s%d_%d", constant.StrategyAwardCountKey, strategyID, awardID)
}
