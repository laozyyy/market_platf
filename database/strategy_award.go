package database

import (
	"big_market/common"
	"big_market/model"
	"gorm.io/gorm"
)

func QueryStrategyAwardList(db *gorm.DB) (result []*model.StrategyAward, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("strategy_award").Find(&result).Error
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	return
}

func QueryStrategyAwardListByStrategyId(db *gorm.DB, strategyId int64) (result []*model.StrategyAward, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("strategy_award").Where("strategy_id = ?", strategyId).Find(&result).Error
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	return
}
