package cache

import (
	"big_market/common/log"
	"big_market/conf"
	"github.com/redis/go-redis/v9"
)

var (
	Client *redis.Client
	url    string
)

func init() {
	config := conf.LoadConfig()
	url = config.Redis.URL
}

func Init() {
	opt, err := redis.ParseURL(url)
	if err != nil {
		log.Errorf("err: %v", err)
		panic(err)
	}

	Client = redis.NewClient(opt)
}
