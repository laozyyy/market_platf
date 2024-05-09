package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"gorm.io/gorm"
)

func QueryRaffleActivityCountByActivityCountID(db *gorm.DB, activityCountID int64) (result *model.RaffleActivityCount, err error) {
	if db == nil {
		db, err = getDB()
	}
	if err != nil {
		log.Errorf("err: %v", err)
		return nil, err
	}
	var tmp []*model.RaffleActivityCount
	err = db.Table("raffle_activity_count").Where("activity_count_id = ?", activityCountID).Find(&tmp).Error
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
