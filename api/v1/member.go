package v1

import (
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/jwt"
	"github.com/junminhong/member-services-center/pkg/smtp"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var postgresDB = postgresql.GetDB()

type registerReq struct {
	Email       string `form:"email" json:"email" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	RepPassword string `form:"rep-password" json:"rep-password" binding:"required"`
}

type response struct {
	Message string `json:"message"`
}

// @Summary 註冊會員帳號
// @Tags member
// @version 1.0
// @Accept application/json
// @produce application/json
// @param data body registerReq true "註冊資料"
// @Success 200 {object} response "{"message":""}"
// @Router /member/register [post]
func Register(c *gin.Context) {
	req := &registerReq{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "請傳送正確資料",
		})
		return
	}
	if strings.Compare(req.Password, req.RepPassword) != 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "請傳送正確資料",
		})
		return
	}
	memberStruct := &model.Member{}
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
	member := &model.Member{}
	err = postgresDB.Where("email = ?", req.Email).First(&member).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "查無此會員訊息",
		})
		return
	}
	if req.Password == member.Password {
		accessToken := jwt.GenerateAccessToken(member.ID)
		c.JSON(http.StatusOK, gin.H{
			"token":      accessToken,
			"email-auth": member.EmailAuth,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "密碼錯誤",
		})
	}
}

func TokenAuth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	tokenParts := strings.Split(token, "Bearer ")
	log.Println(jwt.VerifyAccessToken(tokenParts[1]))
}
