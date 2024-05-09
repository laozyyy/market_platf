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
