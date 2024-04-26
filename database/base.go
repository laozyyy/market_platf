package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	dsn = "root:123456@tcp(localhost:13306)/big_market?charset=utf8mb4&parseTime=True&loc=Local"
)

func getDB() (*gorm.DB, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}
