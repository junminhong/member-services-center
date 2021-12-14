package main

import (
	"github.com/junminhong/member-services-center/app/services/grpc"
	"github.com/junminhong/member-services-center/config/routes"
	"sync"
)

func main() {
	// 專案初始化要用這個建立資料庫
	//db.MigrateDB(database.GetDB())

	// 這裡切換api version是指切換公開出去的api url會隨著版本改變
	// ex: 切換v1 host/api/v1/xxx
	//     切換v2 host/api/v2/xxx
	// 更改api version後還需至routes更改對應的controllers並import
	intiServerWg := &sync.WaitGroup{}
	intiServerWg.Add(2)
	go grpc.InitGRpcServer(intiServerWg)
	go routes.InitRoutes("v1", intiServerWg)
	intiServerWg.Wait()
}
