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
		c.Writer.Header().Set("Content-Type", "application/json")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func Init(apiVersion string, intiServerWg *sync.WaitGroup) *gin.Engine {
	defer intiServerWg.Done()
	router := gin.Default()
	router.Use(CORSMiddleware())
	var apiRouter *gin.RouterGroup
	switch apiVersion {
	case "v1", "V1":
		apiRouter = router.Group("/api/v1")
	}
	memberRouter := apiRouter.Group("/member")
	{
		memberRouter.POST("/register", v1.Register)
		memberRouter.POST("/login", v1.Login)
		memberRouter.PUT("/reset-password", v1.ResetPassword)
		memberRouter.POST("/resend-email", v1.ResendEmail)
		memberRouter.PUT("/profile", v1.EditProfile)
		memberRouter.GET("/profile", v1.GetProfile)
	}
	authRouter := apiRouter.Group("/auth")
	{
		authRouter.GET("/email", v1.VerifyEmail)
		authRouter.GET("/member", v1.TokenAuth)
		authRouter.POST("/resend-email", v1.ResendEmail)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.DefaultModelsExpandDepth(-1)))
	return router
}
