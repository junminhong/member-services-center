package v1

import (
	"context"
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/handler"
	"github.com/junminhong/member-services-center/pkg/jwt"
	"github.com/junminhong/member-services-center/pkg/logger"
	"github.com/junminhong/member-services-center/pkg/smtp"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var postgresDB = postgresql.GetDB()
var sugar = logger.Setup()

type registerReq struct {
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	RepPassword string `json:"rep_password" binding:"required"`
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

type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type loginData struct {
	AccessToken string `json:"access_token"`
}

// Login
// @Summary 登入會員帳號
// @Tags member
// @version 1.0
// @Accept application/json
// @produce application/json
// @param data body LoginReq true "登入資料"
// @Success 200 {object} handler.Response
// @Router /member/login [post]
func Login(c *gin.Context) {
	request := &LoginReq{}
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
	response := loginHandler(request.Email, request.Password)
	c.JSON(response.ResultCode, response)
}
func loginHandler(email string, password string) handler.Response {
	member := &model.Member{}
	err := postgresDB.Where("email = ?", email).First(&member).Error
	var response handler.Response
	if err != nil {
		response = handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "信箱輸入錯誤",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		}
		return response
	}
	if password != member.Password {
		response = handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "密碼輸入錯誤",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		}
		return response
	}
	if !member.EmailAuth {
		response = handler.Response{
			ResultCode: handler.BadRequest,
			Message:    "該用戶信箱未認證",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		}
		return response
	}
	accessToken := jwt.GenerateAccessToken(member.ID)
	response = handler.Response{
		ResultCode: handler.OK,
		Message:    "登入成功",
		Data:       loginData{AccessToken: accessToken},
		TimeStamp:  time.Now().UTC(),
	}
	return response
}

func TokenAuth(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	tokenParts := strings.Split(token, "Bearer ")
	if !jwt.VerifyAccessToken(tokenParts[1]) {
		c.JSON(handler.Forbidden, handler.Response{
			ResultCode: handler.Forbidden,
			Message:    "access token驗證失敗",
			Data:       "",
			TimeStamp:  time.Now().UTC(),
		})
		return
	}
	// 如果token是合法的就拿token去跟redis換member id回來
	memberID, err := redisClient.Get(context.Background(), tokenParts[1]).Result()
	if err != nil {
		sugar.Info(err.Error())
	}
	log.Println(memberID)
}
