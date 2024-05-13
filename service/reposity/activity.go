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
	"time"
)

func GetActivityByActivityID(activityID int64) (*model.RaffleActivity, error) {
	key := fmt.Sprintf("%s%d", constant.ActivityKey, activityID)
	result, err := cache.Client.Get(context.Background(), key).Result()
	if result != "" {
		var activity model.RaffleActivity
		err = json.Unmarshal([]byte(result), &activity)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		return &activity, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("err: %v", err)
		return nil, err
	}
	activity, err := database.QueryRaffleActivityByActivityID(nil, activityID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	marshal, err := json.Marshal(activity)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = cache.Client.Set(context.Background(), key, marshal, 0).Err()
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return activity, nil
}

func GetActivityCountByActivityCountID(activityCountID int64) (*model.RaffleActivityCount, error) {
	key := fmt.Sprintf("%s%d", constant.ActivityCountKey, activityCountID)
	result, err := cache.Client.Get(context.Background(), key).Result()
	if result != "" {
		var activityCount model.RaffleActivityCount
		err = json.Unmarshal([]byte(result), &activityCount)
		if err != nil {
			log.Errorf("err: %v", err)
			return nil, err
		}
		return &activityCount, nil
	}
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("err: %v", err)
		return nil, err
	}
	activityCount, err := database.QueryRaffleActivityCountByActivityCountID(nil, activityCountID)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	marshal, err := json.Marshal(activityCount)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = cache.Client.Set(context.Background(), key, marshal, 0).Err()
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return activityCount, nil
}

func CacheSku(skuEntity model.RaffleActivitySku) error {
	key := fmt.Sprintf("%s%d", constant.ActivitySkuStockCountKey, skuEntity.SKU)
	err := cache.Client.Set(context.Background(), key, skuEntity.StockCountSurplus, 0).Err()
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

func DescActivitySkuStock(sku int64, endTime time.Time) (int, error) {
	key := fmt.Sprintf("%s%d", constant.ActivitySkuStockCountKey, sku)
	surplus, err := cache.Client.Decr(context.Background(), key).Result()
	if err != nil {
		log.Errorf("err: %v", err)
		return 0, err
	}
	// 恢复库存
	if surplus < 0 {
		log.Infof("库存为0 sku: %d", sku)
		err = cache.Client.Set(context.Background(), key, 0, 0).Err()
		if err != nil {
			log.Errorf("error: %v", err)
			return 0, err
		}
	}
	// setnx
	lockKey := fmt.Sprintf("%s_%d", key, surplus)
	//expireTime := endTime.AddDate(0, 0, 1).Sub(time.Now())
	expireTime := time.Second
	result, err := cache.Client.SetNX(context.Background(), lockKey, 1, expireTime).Result()
	if err != nil {
		log.Errorf("error: %v", err)
		return 0, err
	}
	if !result {
		log.Errorf("加锁失败, surplus: %d", surplus)
		return 0, nil
	}
	if surplus == 0 {
		return 2, nil
	}
	return 1, nil
}
