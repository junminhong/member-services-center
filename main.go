package main

import (
	"github.com/junminhong/member-services-center/app/api/v1/routes"
)

func main() {
	// 這裡切換api version是指切換公開出去的api url會隨著版本改變
	// ex: 切換v1 host/api/v1/xxx
	//     切換v2 host/api/v2/xxx
	// 更改api version後還需至routes更改對應的controllers並import
	routes.InitRoutes("v1")
}
