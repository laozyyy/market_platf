package database

import (
	"big_market/common/log"
	"big_market/model"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/sharding"
	"hash/fnv"
)

func getShardingMiddleware() *sharding.Sharding {
	middleware := sharding.Register(sharding.Config{
		ShardingKey:    "user_id",
		NumberOfShards: 4,
		ShardingAlgorithm: func(columnValue any) (suffix string, err error) {
			if uid, ok := columnValue.(string); ok {
				if uid == "xiaofuge" {
					return "_001", nil
				}
				h := fnv.New32a()
				h.Write([]byte(uid))
				parseInt := h.Sum32()
				if err != nil {
					return "", errors.New("invalid user_id")
				}
				suffix = fmt.Sprintf("_00%d", parseInt%4)
				log.Infof("表后缀%s", suffix)
				return suffix, nil
			}
			return "", errors.New("invalid user_id")
		},
		PrimaryKeyGenerator: sharding.PKCustom,
		PrimaryKeyGeneratorFn: func(tableIdx int64) int64 {
			return 0
		},
	}, "raffle_activity_order")
	return middleware
}

// InsertRaffleActivityOrder 分片
func InsertRaffleActivityOrder(db *gorm.DB, order *model.RaffleActivityOrder) (err error) {
	if db == nil {
		db, err = getDB()
	}
	db.Use(getShardingMiddleware())
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	err = db.Table("raffle_activity_order").Omit("id").Create(order).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return
}

func UpdateRaffleActivityOrderTotalCountByUserID(db *gorm.DB, uID string, count int) (err error) {
	if db == nil {
		db, err = getDB()
	}
	db.Use(getShardingMiddleware())
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	err = db.Table("raffle_activity_order").
		Where("user_id", uID).
		Update("total_count", count).Error
	if err != nil {
		log.Errorf("err: %v", err)
		return err
	}
	return
}
