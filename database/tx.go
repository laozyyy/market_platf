package database

import (
	"big_market/model"
	"gorm.io/gorm"
)

func SaveOrderAndAccount(order model.RaffleActivityOrder, account model.RaffleActivityAccount) error {
	db := DB
	err := db.Transaction(func(tx *gorm.DB) error {
		// 插入订单
		if err := InsertRaffleActivityOrder(tx, &order); err != nil {
			// 返回任何错误都会回滚事务
			return err
		}
		affected, err := UpdateRaffleActivityAccount(tx, &account)
		if err != nil {
			return err
		}
		if affected == 0 {
			if err = InsertRaffleActivityAccount(tx, &account); err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
		}
		// 提交事务
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
