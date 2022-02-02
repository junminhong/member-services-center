package v1

import (
	"context"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/handler"
	"time"

	"github.com/gin-gonic/gin"
)

var redisClient = redis.Setup()

func updateEmailAuth(email string) bool {
	err := postgresDB.Model(&model.Member{}).Where("email = ?", email).Update("email_auth", true).Error
	if err != nil {
		sugar.Info(err.Error())
	}
	return err != nil
}

func VerifyEmail(c *gin.Context) {
	email, err := redisClient.Get(context.Background(), c.Query("email_token")).Result()
	if err != nil {
		sugar.Info(err.Error())
	}
	if email == "" {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError2,
			Message:    handler.ResponseFlag[handler.AuthError2],
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	if updateEmailAuth(email) {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError3,
			Message:    handler.ResponseFlag[handler.AuthError3],
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	err = redisClient.Del(context.Background(), c.Param("emailToken")).Err()
	if err != nil {
		sugar.Info(err.Error())
	}
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.AuthOK1,
		Message:    handler.ResponseFlag[handler.AuthOK1],
		Data:       "",
		TimeStamp:  time.Now().UTC(),
	})
}
