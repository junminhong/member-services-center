package router

import (
	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/api/v1"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

func SetupRouter(apiVersion string, intiServerWg *sync.WaitGroup) *gin.Engine {
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
		memberRouter.POST("/register", v1.Register)
		memberRouter.POST("/login", v1.Login)
		memberRouter.POST("/email-auth", v1.VerifyEmail)
		memberRouter.POST("/token-auth", v1.TokenAuth)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	return router
}