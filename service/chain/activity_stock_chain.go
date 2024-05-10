package chain

import (
	"big_market/common/log"
	"big_market/model"
)

type ActivityStockChain struct {
	nextChain ActivityChain
}

func (a *ActivityStockChain) AppendNext(next *ActivityChain) *ActivityChain {
	a.nextChain = *next
	return next
}

func (a *ActivityStockChain) Action(activity model.RaffleActivity, sku model.RaffleActivitySku, count model.RaffleActivityCount) error {
	log.Infof("责任链：活动库存处理, activyty: %v sku: %v, count: %v", activity, sku, count)

	return nil
}
