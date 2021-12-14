package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/app/api/v1/controllers/member"
)

func InitRoutes(apiVersion string) {
	router := gin.Default()
	var apiVersionTmp *gin.RouterGroup
	switch apiVersion {
	case "v1", "V1":
		apiVersionTmp = router.Group("/api/v1")
	}
	memberRouter := apiVersionTmp.Group("/member")
	{
		memberRouter.POST("/register", member.Register)
		memberRouter.POST("/login", member.Login)
		memberRouter.POST("/email-auth", member.VerifyEmail)
	}
	router.Run(":8080")
}
