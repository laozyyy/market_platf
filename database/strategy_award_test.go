package database

import (
	"big_market/redis"
	"context"
	"fmt"
	"testing"
)

func Test_queryStrategyAwardList(t *testing.T) {
	ctx := context.Background()
	result, _ := redis.Client.HGetAll(ctx, "strategy_id_100001").Result()
	fmt.Println(result)
}
