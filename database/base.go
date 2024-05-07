package database

import (
	"big_market/conf"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dsn string
)

// DB 外部调用的单例
var DB *gorm.DB

func getDB() (*gorm.DB, error) {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return nil, err
	}
	return db, nil
}

func init() {
	config := conf.LoadConfig()
	dsn = config.Database.URL
}

func Init() {
	var err error
	DB, err = getDB()
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
	}
}
