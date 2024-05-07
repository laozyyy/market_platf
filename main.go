package main

import (
	"big_market/cache"
	"big_market/crons"
	"big_market/database"
	"big_market/mq"
	"github.com/gin-gonic/gin"
)

func main() {

	cache.Init()
	database.Init()
	mq.Init()
	crons.AddCron()
	r := gin.Default()
	initRouter(r)
	_ = r.Run()
	//dsn := "postgres://localhost:5432/sharding-db?sslmode=disable"
	//db, err := gorm.Open(mysql.New(mysql.Config{DSN: dsn}))
	//if err != nil {
	//	panic(err)
	//}
	//
	//for i := 0; i < 64; i += 1 {
	//	table := fmt.Sprintf("orders_%02d", i)
	//	db.Exec(`DROP TABLE IF EXISTS ` + table)
	//	db.Exec(`CREATE TABLE ` + table + ` (
	//		id BIGSERIAL PRIMARY KEY,
	//		user_id bigint,
	//		product_id bigint
	//	)`)
	//}
	//
	//middleware := sharding.Register(sharding.Config{
	//	ShardingKey:         "user_id",
	//	NumberOfShards:      64,
	//	PrimaryKeyGenerator: sharding.PKSnowflake,
	//}, "orders")
	//db.Use(middleware)

}
