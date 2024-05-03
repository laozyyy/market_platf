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

func GetStrategyByStrategyID(strategyID int64) (strategy *model.Strategy, err error) {
	ctx := context.Background()
	strategyStr, err := cache.Client.Get(ctx, fmt.Sprintf("%s%d", constant.StrategyKey, strategyID)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Errorf("error: %v", err)
		return nil, err
	}
	if strategyStr != "" {
		err = json.Unmarshal([]byte(strategyStr), &strategy)
		if err != nil {
			log.Errorf("error: %v", err)
			return nil, err
		}
		return
	}
	strategy, err = database.QueryStrategyByStrategyID(nil, strategyID)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}
	//缓存
	marshal, err := json.Marshal(strategy)
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}
	err = cache.Client.Set(ctx, fmt.Sprintf("%s%d", constant.StrategyKey, strategyID), marshal, 0).Err()
	if err != nil {
		log.Errorf("error: %v", err)
		return nil, err
	}
	return
}
