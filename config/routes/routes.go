package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/app/api/v1/controllers/member"
	"sync"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func InitRoutes(apiVersion string, intiServerWg *sync.WaitGroup) {
	defer intiServerWg.Done()
	router := gin.Default()
	router.Use(CORSMiddleware())
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
