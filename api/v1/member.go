package v1

import (
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/handler"
	"github.com/junminhong/member-services-center/pkg/jwt"
	"github.com/junminhong/member-services-center/pkg/logger"
	"github.com/junminhong/member-services-center/pkg/smtp"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var postgresDB = postgresql.GetDB()
var sugar = logger.Setup()

type registerReq struct {
	Email       string `form:"email" json:"email" binding:"required"`
	Password    string `form:"password" json:"password" binding:"required"`
	RepPassword string `form:"rep-password" json:"rep-password" binding:"required"`
}

// Register
// @Summary 註冊會員帳號
// @Tags member
// @version 1.0
// @Accept application/json
// @produce application/json
// @param data body registerReq true "註冊資料"
// @Success 200 {object} handler.Response
// @Router /member/register [post]
func Register(c *gin.Context) {
	request := &registerReq{}
	err := c.BindJSON(request)
	if err != nil {
		c.JSON(handler.BadRequest, handler.Response{
			ResultCode: handler.BadRequest,
			Message:    handler.ResponseFlag[handler.BadRequest],
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	if strings.Compare(request.Password, request.RepPassword) != 0 {
		c.JSON(handler.BadRequest, handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "密碼不匹配，請重新請求",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	response := registerHandler(request)
	if response.ResultCode != handler.OK {
		c.JSON(response.ResultCode, response)
		return
	}
	smtp.SendEmailAuth(request.Email)
	c.JSON(response.ResultCode, response)
}

func registerHandler(request *registerReq) *handler.Response {
	memberStruct := &model.Member{}
	memberStruct.CreatedAt = time.Now().UTC()
	memberStruct.UpdatedAt = time.Now().UTC()
	memberStruct.ActivatedAt = time.Now().UTC()
	memberStruct.Email = request.Email
	memberStruct.Password = request.Password
	result := postgresDB.Create(memberStruct)
	var response *handler.Response
	if result.Error == nil {
		response = &handler.Response{
			ResultCode: handler.OK,
			Message:    "帳戶註冊成功",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		}
	} else {
		response = &handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "帳戶註冊失敗",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		}
		sugar.Info(result.Error)
	}
	return response
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
