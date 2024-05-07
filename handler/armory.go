package handler

import (
	"big_market/common/log"
	"big_market/model/response"
	"big_market/service/reposity"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func HandleArmory(ctx *gin.Context) {
	var params map[string]interface{}
	err := ctx.ShouldBindJSON(&params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "参数绑定失败")
	}
	strategyIDs, ok := params["strategy_ids"].([]interface{})
	if !ok {
		ctx.JSON(http.StatusInternalServerError, "参数解析失败")
	}
	for _, strategyID := range strategyIDs {
		number, ok := strategyID.(float64)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, "参数解析失败")
		}
		err = reposity.AssembleLotteryStrategyWithRules(int64(number))
		if err != nil {
			log.Infof("装配失败")
			ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("装配失败strategyID: %d, err: %v", int64(number), err))
		} else {
			log.Infof("装配成功, strategyID: %d", int64(number))
		}
	}
	ctx.JSON(http.StatusOK, "success")
}

func HandleQueryStrategyAwardList(ctx *gin.Context) {
	strategyID := ctx.Param("strategy_id")
	parseInt, err := strconv.ParseInt(strategyID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "参数解析错误")
	}
	strategyAwardList, err := reposity.GetStrategyAwardList(parseInt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "服务内部错误")
	}
	strategyAwardDTOS := make([]response.StrategyAwardDTO, 0)
	for _, strategyAward := range strategyAwardList {
		dto := response.StrategyAwardDTO{
			AwardId:       strategyAward.AwardId,
			AwardTitle:    strategyAward.AwardTitle,
			AwardSubtitle: strategyAward.AwardSubtitle,
			Sort:          strategyAward.Sort,
		}
		strategyAwardDTOS = append(strategyAwardDTOS, dto)
	}
	resp := response.StrategyAwardListResponse{
		Data: response.StrategyAwardPageInfo{
			//其他字段先不加
			StrategyAward: strategyAwardDTOS,
		},
		Message: "success",
	}
	ctx.JSON(http.StatusOK, resp)
}
