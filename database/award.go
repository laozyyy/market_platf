package database

import (
	"big_market/common"
	"big_market/model"
	"gorm.io/gorm"
)

func QueryAwardList(db *gorm.DB) (result []*model.Award, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("award").Find(&result).Error
	if err != nil {
		common.Log.Errorf("err: %v", err)
		return nil, err
	}
	return
}
