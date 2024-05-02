package cache

import (
	"big_market/common/log"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
)

func init() {
	opt, err := redis.ParseURL("redis://root:@localhost:16379/0")
	if err != nil {
		log.Errorf("err: %v", err)
		panic(err)
	}

	Client = redis.NewClient(opt)
}
