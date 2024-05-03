package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryRuleTreeNodesByTreeID(db *gorm.DB, treeID string) (result []model.RuleTreeNode, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("rule_tree_node").Where("tree_id = ?", treeID).Find(&result).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	if len(result) > 0 {
		return
	} else {
		return nil, errors.New("no data")
	}
}
