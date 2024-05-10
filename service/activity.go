package service

import (
	"big_market/common"
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/database"
	"big_market/model"
	"big_market/service/chain"
	"big_market/service/reposity"
	"time"
)

func CreateRaffleActivityOrder(order model.SkuRechargeEntity) (string, error) {
	sku, err := database.QueryRaffleActivitySkuBySku(nil, order.Sku)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	activity, err := reposity.GetActivityByActivityID(sku.ActivityID)
	log.Infof("activity: %v", activity)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	count, err := reposity.GetActivityCountByActivityCountID(sku.ActivityCountID)
	activityChain := openActivityChain()
	err = activityChain.Action(*activity, *sku, *count)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	aggregate := buildOrderAggregate(order, *activity, *sku, *count)
	err = doSaveOrder(aggregate)
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	return aggregate.ActivityOrderEntity.OrderID, nil
}

func doSaveOrder(aggregate model.CreateOrderAggregate) error {
	order := model.RaffleActivityOrder{
		UserID:        aggregate.UserID,
		SKU:           aggregate.ActivityOrderEntity.Sku,
		ActivityID:    aggregate.ActivityID,
		ActivityName:  aggregate.ActivityOrderEntity.ActivityName,
		StrategyID:    aggregate.ActivityOrderEntity.ActivityID,
		OrderID:       aggregate.ActivityOrderEntity.OrderID,
		OrderTime:     aggregate.ActivityOrderEntity.OrderTime,
		TotalCount:    aggregate.TotalCount,
		DayCount:      aggregate.DayCount,
		MonthCount:    aggregate.MonthCount,
		State:         aggregate.ActivityOrderEntity.State,
		OutBusinessNo: aggregate.ActivityOrderEntity.OutBusinessNo,
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	}
	account := model.RaffleActivityAccount{
		UserID:            aggregate.UserID,
		ActivityID:        aggregate.ActivityID,
		TotalCount:        aggregate.TotalCount,
		TotalCountSurplus: aggregate.TotalCount,
		DayCount:          aggregate.DayCount,
		DayCountSurplus:   aggregate.DayCount,
		MonthCount:        aggregate.MonthCount,
		MonthCountSurplus: aggregate.MonthCount,
		CreateTime:        time.Now(),
		UpdateTime:        time.Now(),
	}
	err := database.SaveOrderAndAccount(order, account)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}

func buildOrderAggregate(order model.SkuRechargeEntity, activity model.RaffleActivity, sku model.RaffleActivitySku, count model.RaffleActivityCount) model.CreateOrderAggregate {
	aggregate := model.CreateOrderAggregate{
		UserID:     order.UserID,
		ActivityID: activity.ActivityID,
		TotalCount: count.TotalCount,
		DayCount:   count.DayCount,
		MonthCount: count.MonthCount,
		ActivityOrderEntity: model.ActivityOrderEntity{
			UserID:        order.UserID,
			ActivityID:    activity.ActivityID,
			ActivityName:  activity.ActivityName,
			StrategyID:    activity.StrategyID,
			OrderID:       common.GenerateRandomString(12),
			OrderTime:     time.Now(),
			TotalCount:    count.TotalCount,
			DayCount:      count.DayCount,
			MonthCount:    count.MonthCount,
			State:         activity.State,
			Sku:           sku.SKU,
			OutBusinessNo: order.OutBusinessNo,
		},
	}
	return aggregate
}

func openActivityChain() chain.ActivityChain {
	activityChain := chain.ActivityChainGroup[constant.ActivityBase]
	next := chain.ActivityChainGroup[constant.ActivityStock]
	activityChain.AppendNext(&next)
	return activityChain
}
