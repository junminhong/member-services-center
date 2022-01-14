package v1

import (
	"context"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var ctx = context.Background()

type emailAuthReq struct {
	EmailToken string `form:"emailToken" json:"emailToken" binding:"required"`
}

func updateEmailAuth(email string) bool {
	err := postgresDB.Model(&model.Member{}).Where("email = ?", email).Update("email_auth", true).Error
	return err != nil
}

func VerifyEmail(c *gin.Context) {
	req := &emailAuthReq{}
	err := c.BindJSON(req)
	if err != nil {
		log.Println(err.Error())
	}
	redisClient := redis.InitRedis()
	email, err := redisClient.Get(ctx, req.EmailToken).Result()
	if err != nil {
		log.Println(err.Error())
	}
	if email != "" {
		if !updateEmailAuth(email) {
			c.JSON(http.StatusOK, gin.H{
				"result":  "200",
				"message": "認證成功",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"result":  "400",
				"message": "完蛋～資料更新失敗！",
			})
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"result":  "400",
			"message": "你是誰？請不要亂來驗證信箱好嗎！",
		})
		return
	}

}
