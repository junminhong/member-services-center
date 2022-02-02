package db

import (
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/logger"
)

var postgresDB = postgresql.GetDB()
var sugar = logger.Setup()

func MigrateDB() {
	defer sugar.Sync()
	sugar.Info("migration db...")
	err := postgresDB.AutoMigrate(&model.Member{}, &model.MemberInfo{})
	if err != nil {
		sugar.Error(err.Error())
	}
}
