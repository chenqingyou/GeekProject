package ioc

import (
	"GeekProject/newGeekProject/day2/webook/config"
	"GeekProject/newGeekProject/day2/webook/internal/repository/dao"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (db *gorm.DB) {
	//初始化数据库
	db, err := gorm.Open(mysql.Open(config.Config.DB.DNS))
	if err != nil {
		fmt.Printf("Open DB init err [%v]\n", err)
		panic(err)
	}
	err = dao.InitDBTable(db) //创建表
	if err != nil {
		fmt.Printf("DB InitDBTable err [%v]\n", err)
		panic(err)
	}
	return db
}
