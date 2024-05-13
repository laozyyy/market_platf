package chain

import (
	"big_market/common/log"
	"big_market/model"
	"big_market/mq"
	"big_market/service/reposity"
	"errors"
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
	flag, err := reposity.DescActivitySkuStock(sku.SKU, activity.EndDateTime)
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	if flag != 0 {
		log.Infof("库存扣减成功, sku:%v", sku)
		if flag == 2 {
			// 库存为0，立即发消息
			err = mq.SendSkuCountZeroMessage(sku.SKU)
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
		} else {
			// 库存非0，延迟消息
			err = mq.SendUpdateSkuCountMessage(sku.SKU)
			if err != nil {
				log.Errorf("err: %v", err)
				return err
			}
		}
		return nil
	}
	log.Errorf("库存扣减失败, sku:%v", sku)
	return errors.New("stock subtraction failed")
}
