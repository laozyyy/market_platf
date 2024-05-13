package chain

import (
	"big_market/common/constant"
	"big_market/common/log"
	"big_market/model"
	"errors"
	"time"
)

type ActivityBaseChain struct {
	nextChain ActivityChain
}

func (a *ActivityBaseChain) AppendNext(next *ActivityChain) *ActivityChain {
	a.nextChain = *next
	return next
}

func (a *ActivityBaseChain) Action(activity model.RaffleActivity, sku model.RaffleActivitySku, count model.RaffleActivityCount) error {
	log.Infof("活动责任链：基础信息校验, activity: %v sku: %v, count: %v", activity, sku, count)
	if activity.State != constant.ActivityOpen {
		return errors.New("activity doesnt open")
	}
	if activity.BeginDateTime.After(time.Now()) ||
		activity.EndDateTime.Before(time.Now()) {
		return errors.New("activity has expired")
	}
	if sku.StockCountSurplus <= 0 {
		return errors.New("activity out of stock")
	}
	err := a.nextChain.Action(activity, sku, count)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
