package database

import (
	"big_market/common/log"
	"big_market/model"
	"gorm.io/gorm"
)

func QueryAwardList(db *gorm.DB) (result []*model.Award, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	err = db.Table("award").Find(&result).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	return
}
