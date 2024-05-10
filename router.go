package main

import (
	"big_market/handler"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	r.POST("/armory", handler.HandleArmory)
	r.POST("/raffle", handler.HandleRaffle)
	r.GET("/query/strategy_award/:strategy_id", handler.HandleQueryStrategyAwardList)
	r.POST("/create/order", handler.HandleCreateOrder)
}
