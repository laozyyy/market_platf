package chain

import (
	"big_market/common/log"
	"big_market/model"
)

type ActivityBaseChain struct {
	nextChain ActivityChain
}

func (a *ActivityBaseChain) AppendNext(next *ActivityChain) *ActivityChain {
	a.nextChain = *next
	return next
}

func (a *ActivityBaseChain) Action(activity model.RaffleActivity, sku model.RaffleActivitySku, count model.RaffleActivityCount) error {
	log.Infof("责任链：活动默认处理, activyty: %v sku: %v, count: %v", activity, sku, count)
	err := a.nextChain.Action(activity, sku, count)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return nil
}
