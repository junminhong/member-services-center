package db

import (
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"log"

	"gorm.io/gorm"
)

var PostgresDB = postgresql.GetDB()

func MigrateDB(db *gorm.DB) {
	log.Println("初始化DB Data")
	db.AutoMigrate(&model.Member{})
}
