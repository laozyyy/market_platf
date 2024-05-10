package handler

import (
	"big_market/model"
	"big_market/model/request"
	"big_market/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func HandleCreateOrder(ctx *gin.Context) {
	var param request.CreateOrderRequest
	err := ctx.BindJSON(&param)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "参数绑定失败")
	}
	skuRechargeEntity := model.SkuRechargeEntity{
		UserID:        param.UserID,
		Sku:           param.Sku,
		OutBusinessNo: param.OutBusinessNo,
	}
	orderID, err := service.CreateRaffleActivityOrder(skuRechargeEntity)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("服务内部错误：%v", err))
	}
	ctx.JSON(http.StatusOK, orderID)
}
