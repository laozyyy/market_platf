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
)

func getStrategyByStrategyID(strategyID int64) (strategy *model.Strategy, err error) {
	ctx := context.Background()
	strategyStr, err := cache.Client.Get(ctx, fmt.Sprintf("%s%d", model.StrategyKey, strategyID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		common.Log.Errorf("error: %v", err)
		return nil, err
	}
	if strategyStr != "" {
		err = json.Unmarshal([]byte(strategyStr), &strategy)
		if err != nil {
			common.Log.Errorf("error: %v", err)
			return nil, err
		}
		return
	}
	strategy, err = database.QueryStrategyByStrategyID(nil, strategyID)
	if err != nil {
		common.Log.Errorf("error: %v", err)
		return nil, err
	}
	//缓存
	marshal, err := json.Marshal(strategy)
	if err != nil {
		common.Log.Errorf("error: %v", err)
		return nil, err
	}
	err = cache.Client.Set(ctx, fmt.Sprintf("%s%d", model.StrategyKey, strategyID), marshal, 0).Err()
	if err != nil {
		common.Log.Errorf("error: %v", err)
		return nil, err
	}
	return
}
