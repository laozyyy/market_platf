package database

import (
	"big_market/model"
	"gorm.io/gorm"
)

func UpdateRaffleActivityAccount(db *gorm.DB, account *model.RaffleActivityAccount) (affected int, err error) {
	if db == nil {
		db, err = getDB()
	}
	result := db.Model(account).
		Table("raffle_activity_account").
		Where("user_id = ?", account.UserID).
		Where("activity_id = ?", account.ActivityID).
		Omit("create_time").
		Updates(account)
	return int(result.RowsAffected), result.Error
}

func InsertRaffleActivityAccount(db *gorm.DB, account *model.RaffleActivityAccount) (err error) {
	if db == nil {
		db, err = getDB()
	}
	err = db.Table("raffle_activity_account").Create(account).Error
	return err
}
