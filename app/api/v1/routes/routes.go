package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/app/api/v1/controllers/member"
)

func InitRoutes() {
	router := gin.Default()
	v1 := router.Group("/api/v1")
	memberRouter := v1.Group("/member")
	{
		memberRouter.POST("/register", member.Register)
	}
	router.Run(":8080")
}
