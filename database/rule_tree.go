package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryRuleTreeByTreeID(db *gorm.DB, treeID string) (result *model.RuleTree, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	var tmp []*model.RuleTree
	err = db.Table("rule_tree").Where("tree_id = ?", treeID).Find(&tmp).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	if len(tmp) > 0 {
		result = tmp[0]
	} else {
		return nil, errors.New("no data")
	}
	return
}
