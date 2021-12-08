package member

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	fmt.Println("收到註冊請求")
}
