package db

import (
	"log"

	"github.com/junminhong/member-services-center/app/api/v1/models/member"
	"github.com/junminhong/member-services-center/config/database"
	"gorm.io/gorm"
)

var PostgresDB = database.GetDB()

func MigrateDB(db *gorm.DB) {
	log.Println("初始化DB Data")
	db.AutoMigrate(&member.Member{})
}
