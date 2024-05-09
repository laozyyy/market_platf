package main

import (
	"big_market/cache"
	"big_market/crons"
	"big_market/database"
	"big_market/mq"
	"github.com/gin-gonic/gin"
)

func main() {
	cache.Init()
	database.Init()
	mq.Init()
	crons.AddCron()
	r := gin.Default()
	initRouter(r)
	_ = r.Run()
}
