package db

import (
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/logger"
)

var postgresDB = postgresql.GetDB()
var sugar = logger.Init()

func MigrateDB() {
	defer sugar.Sync()
	if !postgresDB.Migrator().HasTable(&model.Member{}) {
		return
	}
	sugar.Info("migration db...")
	err := postgresDB.AutoMigrate(&model.Member{})
	if err != nil {
		sugar.Error(err.Error())
	}
}
