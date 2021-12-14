package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/app/api/v1/controllers/member"
	"sync"
)

func InitRoutes(apiVersion string, intiServerWg *sync.WaitGroup) {
	defer intiServerWg.Done()
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
		memberRouter.POST("/token-auth", member.TokenAuth)
	}
	router.Run(":8080")
}
