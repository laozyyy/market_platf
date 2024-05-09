package service

import (
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"big_market/service/reposity"
	"time"
)

func CreateRaffleActivityOrder(order model.ActivityShopCartEntity) (*model.ActivityOrderEntity, error) {
	sku, err := database.QueryRaffleActivitySkuBySku(nil, order.Sku)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	activity, err := reposity.GetActivityByActivityID(sku.ActivityID)
	log.Infof("activity: %v", activity)
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	count, err := reposity.GetActivityCountByActivityCountID(sku.ActivityCountID)
	return &model.ActivityOrderEntity{
		UserID:       order.UserID,
		ActivityID:   sku.ActivityID,
		ActivityName: activity.ActivityName,
		StrategyID:   0,
		OrderID:      "",
		OrderTime:    time.Now(),
		TotalCount:   count.TotalCount,
		DayCount:     count.DayCount,
		MonthCount:   count.MonthCount,
		State:        activity.State,
	}, nil
}
