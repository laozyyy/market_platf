package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryRaffleActivityByActivityID(db *gorm.DB, activityID int64) (result *model.RaffleActivity, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	var tmp []*model.RaffleActivity
	err = db.Table("raffle_activity").Where("activity_id = ?", activityID).Find(&tmp).Error
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
