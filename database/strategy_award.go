package database

import (
	"big_market/common"
	"big_market/common/log"
	"big_market/model"
	"gorm.io/gorm"
)

func QueryStrategyAwardList(db *gorm.DB) (result []*model.StrategyAward, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("strategy_award").Find(&result).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return
}

func QueryStrategyAwardListByStrategyId(db *gorm.DB, strategyID int64) (result []*model.StrategyAward, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("strategy_award").Where("strategy_id = ?", strategyID).Find(&result).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return
}

func QueryStrategyAwardRuleModel(db *gorm.DB, strategyID int64, awardID int) (result string, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	var tmp []model.StrategyAward
	err = db.Table("strategy_award").
		Where("award_id = ?", awardID).
		Where("strategy_id = ?", strategyID).
		Find(&tmp).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return "", err
	}
	if len(tmp) > 0 {
		result = tmp[0].RuleModels
		return
	} else {
		return "", common.NoDataErr
	}
}
