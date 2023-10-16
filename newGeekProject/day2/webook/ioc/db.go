package ioc

import (
	"GeekProject/newGeekProject/day2/webook/internal/repository/dao"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() (db *gorm.DB) {
	type DBConfig struct {
		Dns string `json:"dns"`
	}
	var cfg DBConfig
	err := viper.UnmarshalKey("db.mysql.dsn", &cfg)
	if err != nil {
		panic(err)
	}
	//初始化数据库
	db, err = gorm.Open(mysql.Open(cfg.Dns))
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
