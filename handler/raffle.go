package handler

import (
	"big_market/model"
	"big_market/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleRaffle(ctx *gin.Context) {
	var param model.RaffleFactor
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusInternalServerError, "参数解析失败")
		return
	}
	result, err := service.PerformRaffle(param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "抽奖失败")
	}
	ctx.JSON(http.StatusOK, result)
}
