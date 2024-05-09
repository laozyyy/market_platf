package main

import (
	"big_market/database"
	"big_market/model"
	"testing"
	"time"
)

func TestSharding(t *testing.T) {
	//dsn := "root:123456@tcp(localhost:13306)/big_market?charset=utf8mb4&parseTime=True&loc=Local"
	//db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}))
	//if err != nil {
	//	panic(err)
	//}
	//num := make([]int, 4)
	//middleware := sharding.Register(sharding.Config{
	//	ShardingKey:    "user_id",
	//	NumberOfShards: 4,
	//	//ShardingSuffixs: func() (suffixs []string) {
	//	//	return []string{"_000", "_001", "_002", "_003"}
	//	//},
	//	ShardingAlgorithm: func(columnValue any) (suffix string, err error) {
	//		if uid, ok := columnValue.(string); ok {
	//			h := fnv.New32a()
	//			h.Write([]byte(uid))
	//			parseInt := h.Sum32()
	//			if err != nil {
	//				return "", errors.New("invalid user_id")
	//			}
	//			suffix = fmt.Sprintf("_00%d", parseInt%4)
	//			num[parseInt%4]++
	//			log.Infof("表后缀%s", suffix)
	//			return suffix, nil
	//		}
	//		return "", errors.New("invalid user_id")
	//	},
	//	PrimaryKeyGenerator: sharding.PKSnowflake,
	//}, "raffle_activity_order")
	//db.Use(middleware)
	_ = database.InsertRaffleActivityOrder(nil, &model.RaffleActivityOrder{
		//ID:            0,
		UserID:        GenerateRandomString(10),
		SKU:           0,
		ActivityID:    0,
		ActivityName:  "",
		StrategyID:    0,
		OrderID:       GenerateRandomString(10),
		OrderTime:     time.Now(),
		TotalCount:    0,
		DayCount:      0,
		MonthCount:    0,
		State:         "",
		OutBusinessNo: GenerateRandomString(10),
		CreateTime:    time.Now(),
		UpdateTime:    time.Now(),
	})
	print(1)
	select {}
}
