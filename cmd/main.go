package main

import (
	"github.com/junminhong/member-services-center/db"
	_ "github.com/junminhong/member-services-center/docs"
	"github.com/junminhong/member-services-center/grpc"
	"github.com/junminhong/member-services-center/pkg/logger"
	"github.com/junminhong/member-services-center/router"
	"os"
	"sync"
)

// @title           Member Center Service API
// @version         1.0
// @description     This is a base golang develop member center service
// @termsOfService  http://swagger.io/terms/

// @contact.name   junmin.hong
// @contact.url    https://github.com/junminhong
// @contact.email  junminhong1110@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth
var sugar = logger.Init()

func main() {
	// 專案初始化要用這個建立資料庫
	db.MigrateDB()

	// 這裡切換api version是指切換公開出去的api url會隨著版本改變
	// ex: 切換v1 host/api/v1/xxx
	//     切換v2 host/api/v2/xxx
	// 更改api version後還需至routes更改對應的controllers並import
	defer sugar.Sync()
	setupServerWG := &sync.WaitGroup{}
	setupServerWG.Add(2)
	router := router.Init("v1", setupServerWG)
	go grpc.SetupServer(setupServerWG)
	go router.Run(":" + os.Getenv("HOST_PORT"))
	setupServerWG.Wait()
}
