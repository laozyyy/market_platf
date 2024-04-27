package database

import (
	"big_market/common"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryStrategyRulesByRuleModel(db *gorm.DB, strategyID int64, ruleModel string) (result *model.StrategyRule, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	var tmp []*model.StrategyRule
	err = db.
		Table("strategy_rule").
		Where("strategy_id = ?", strategyID).
		Where("rule_model = ?", ruleModel).
		Find(&tmp).Error
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	if len(tmp) > 0 {
		result = tmp[0]
	} else {
		return nil, errors.New("no data")
	}
	return
}
