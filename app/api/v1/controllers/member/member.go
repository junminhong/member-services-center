package member

import (
	"github.com/junminhong/member-services-center/app/services/jwt"
	"github.com/junminhong/member-services-center/config/database"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/junminhong/member-services-center/app/api/v1/models/member"
	"github.com/junminhong/member-services-center/app/services/smtp"
	"github.com/junminhong/member-services-center/db"
)

var postgresDB = db.PostgresDB

type registerReq struct {
	Email       string `form:"email" json:"email" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	RepPassword string `form:"rep-password" json:"rep-password" binding:"required"`
}

func Register(c *gin.Context) {
	req := &registerReq{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "請傳送正確資料",
		})
		return
	}
	memberStruct := &member.Member{}
	memberStruct.CreatedAt = time.Now().UTC()
	memberStruct.UpdatedAt = time.Now().UTC()
	memberStruct.ActivatedAt = time.Now().UTC()
	memberStruct.Email = req.Email
	memberStruct.Password = req.Password
	result := postgresDB.Create(memberStruct)
	if result.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "註冊成功",
		})
		smtp.SendEmailAuth(req.Email)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "註冊失敗",
		})
	}
}

type loginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	req := &loginReq{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "請傳送正確資料",
		})
		return
	}
	db := database.GetDB()
	member := &member.Member{}
	err = db.Where("email = ?", req.Email).First(&member).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "沒有此位會員",
		})
		return
	}
	accessToken := jwt.GenerateAccessToken(member.ID)
	c.JSON(http.StatusOK, gin.H{
		"message": accessToken,
	})
}

func TokenAuth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	tokenParts := strings.Split(token, "Bearer ")
	log.Println(jwt.VerifyAccessToken(tokenParts[1]))
}
