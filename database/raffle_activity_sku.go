package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryRaffleActivitySkuBySku(db *gorm.DB, sku int64) (result *model.RaffleActivitySku, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	var tmp []*model.RaffleActivitySku
	err = db.Table("raffle_activity_sku").Where("sku = ?", sku).Find(&tmp).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	if len(tmp) > 0 {
		result = tmp[0]
	} else {
		return nil, errors.New("no sku data")
	}
	return
}
