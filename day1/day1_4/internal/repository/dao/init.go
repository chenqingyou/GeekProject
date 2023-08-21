package dao

import "gorm.io/gorm"

func InitDBTable(db *gorm.DB) error {
	return db.AutoMigrate(&UserDB{})
}
