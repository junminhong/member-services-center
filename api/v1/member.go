package v1

import (
	"context"
	"github.com/google/uuid"
	"github.com/junminhong/member-services-center/db/postgresql"
	"github.com/junminhong/member-services-center/model"
	"github.com/junminhong/member-services-center/pkg/gcp"
	"github.com/junminhong/member-services-center/pkg/handler"
	"github.com/junminhong/member-services-center/pkg/jwt"
	"github.com/junminhong/member-services-center/pkg/logger"
	"github.com/junminhong/member-services-center/pkg/smtp"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var postgresDB = postgresql.GetDB()
var sugar = logger.Setup()

type registerReq struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	NickName string `json:"nick_name" binding:"required"`
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
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RequestFormatError1,
			Message:    handler.ResponseFlag[handler.RequestFormatError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	var memberCounts int64
	postgresDB.Where("email = ?", request.Email).Model(&model.Member{}).Count(&memberCounts)
	if memberCounts != 0 {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RegisterError2,
			Message:    handler.ResponseFlag[handler.RegisterError2],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	response := registerHandler(request)
	if response.ResultCode == handler.RegisterError1 {
		c.JSON(handler.OK, response)
		return
	}
	go smtp.SendEmailAuth(request.Email)
	c.JSON(handler.OK, response)
}

func registerHandler(request *registerReq) *handler.Response {
	uuid := uuid.New()
	result := postgresDB.Create(&model.Member{
		UUID:     uuid.String(),
		Email:    request.Email,
		Password: request.Password,
		MemberInfo: model.MemberInfo{
			NickName:    request.NickName,
			MugShotPath: "",
		},
	})
	if result.Error != nil {
		sugar.Info(result.Error)
		return &handler.Response{
			ResultCode: handler.RegisterError1,
			Message:    handler.ResponseFlag[handler.RegisterError1],
			Data:       "",
			TimeStamp:  time.Now(),
		}
	}
	return &handler.Response{
		ResultCode: handler.RegisterOK1,
		Message:    handler.ResponseFlag[handler.RegisterOK1],
		Data:       "",
		TimeStamp:  time.Now(),
	}
}

type LoginReq struct {
	Email    string `form:"email" json:"email" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}
type loginData struct {
	AtomicToken        string `json:"atomic_token"`
	RefreshAtomicToken string `json:"refresh_atomic_token"`
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
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RequestFormatError1,
			Message:    handler.ResponseFlag[handler.RequestFormatError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	response := loginHandler(request.Email, request.Password)
	c.JSON(handler.OK, response)
}

func loginHandler(email string, password string) handler.Response {
	member := &model.Member{}
	err := postgresDB.Where("email = ?", email).First(&member).Error
	var response handler.Response
	if err != nil {
		response = handler.Response{
			ResultCode: handler.LoginError1,
			Message:    handler.ResponseFlag[handler.LoginError1],
			Data:       "",
			TimeStamp:  time.Now(),
		}
		return response
	}
	if strings.Compare(password, member.Password) != 0 {
		response = handler.Response{
			ResultCode: handler.LoginError2,
			Message:    handler.ResponseFlag[handler.LoginError2],
			Data:       "",
			TimeStamp:  time.Now(),
		}
		return response
	}
	if !member.EmailAuth {
		response = handler.Response{
			ResultCode: handler.LoginError3,
			Message:    handler.ResponseFlag[handler.LoginError3],
			Data:       "",
			TimeStamp:  time.Now(),
		}
		return response
	}
	// 登入要給兩個token
	// atomic token 每一小時過期
	// refresh atomic token 每一天過期
	atomicToken := jwt.GenerateAccessToken(member.UUID)
	refreshAtomicToken := jwt.GenerateRefreshAtomicToken(member.UUID)
	response = handler.Response{
		ResultCode: handler.LoginOK1,
		Message:    handler.ResponseFlag[handler.LoginOK1],
		Data: loginData{
			AtomicToken:        atomicToken,
			RefreshAtomicToken: refreshAtomicToken,
		},
		TimeStamp: time.Now(),
	}
	member.AtomicToken = atomicToken
	member.RefreshAtomicToken = refreshAtomicToken
	go postgresDB.Save(member)
	return response
}

func TokenAuth(c *gin.Context) {
	// token := c.Request.Header.Get("Authorization")
	// tokenParts := strings.Split(token, "Bearer ")
	atomicToken := c.Query("atomic_token")
	if authToken(atomicToken) == "" {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError1,
			Message:    handler.ResponseFlag[handler.AuthError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.AuthOK2,
		Message:    handler.ResponseFlag[handler.AuthOK2],
		Data:       "",
		TimeStamp:  time.Now(),
	})
}

func authToken(atomicToken string) string {
	if !jwt.VerifyAtomicToken(atomicToken) {
		return ""
	}
	// 如果token是合法的就拿token去跟redis換member uuid回來
	memberUUID, err := redisClient.Get(context.Background(), atomicToken).Result()
	if err != nil {
		sugar.Info(err.Error())
	}
	return memberUUID
}

type resetPasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func ResetPassword(c *gin.Context) {
	request := &resetPasswordRequest{}
	err := c.BindJSON(request)
	if err != nil {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RequestFormatError1,
			Message:    handler.ResponseFlag[handler.RequestFormatError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	token := c.Request.Header.Get("Authorization")
	tokens := strings.Split(token, " ")
	if len(tokens) != 2 {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError4,
			Message:    handler.ResponseFlag[handler.AuthError4],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	if strings.Compare(tokens[0], "Bearer") != 0 {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError4,
			Message:    handler.ResponseFlag[handler.AuthError4],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	atomicToken := tokens[1]
	memberUUID := authToken(atomicToken)
	if memberUUID == "" {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.AuthError1,
			Message:    handler.ResponseFlag[handler.AuthError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	member := &model.Member{}
	err = postgresDB.Where("uuid = ?", memberUUID).First(&member).Error
	if err != nil {
		sugar.Info(err.Error())
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.ResetPasswordError2,
			Message:    handler.ResponseFlag[handler.ResetPasswordError2],
			Data:       "",
			TimeStamp:  time.Now(),
		})
	}
	if strings.Compare(member.Password, request.OldPassword) != 0 {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.ResetPasswordError1,
			Message:    handler.ResponseFlag[handler.ResetPasswordError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	member.Password = request.NewPassword
	err = postgresDB.Save(&member).Error
	if err != nil {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.ResetPasswordError3,
			Message:    handler.ResponseFlag[handler.ResetPasswordError3],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.ResetPasswordOK1,
		Message:    handler.ResponseFlag[handler.ResetPasswordOK1],
		Data:       "",
		TimeStamp:  time.Now(),
	})
}

type resendEmailRequest struct {
	Email string `json:"email" binding:"required"`
}

func ResendEmail(c *gin.Context) {
	// 每一小時才能寄一次信
	request := &resendEmailRequest{}
	err := c.BindJSON(request)
	if err != nil {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RequestFormatError1,
			Message:    handler.ResponseFlag[handler.RequestFormatError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	if !smtp.SendEmailAuth(request.Email) {
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.SmtpError1,
			Message:    handler.ResponseFlag[handler.SmtpError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.SmtpOK1,
		Message:    handler.ResponseFlag[handler.SmtpOK1],
		Data:       "",
		TimeStamp:  time.Now(),
	})
}

func EditProfile(c *gin.Context) {
	c.JSON(handler.OK, handler.Response{
		Message: "編輯",
	})
}

func GetProfile(c *gin.Context) {
	c.JSON(handler.OK, handler.Response{
		Message: "取得",
	})
}

type uploadMugShotRequest struct {
	MugShot *multipart.FileHeader `form:"mug_shot"`
}

// UploadMugShot 根據token就可以知道該大頭貼是哪一位用戶的
func UploadMugShot(c *gin.Context) {
	request := &uploadMugShotRequest{}
	err := c.Bind(request)
	//file, err := c.FormFile("mug_shot")
	file, uploadFile, err := c.Request.FormFile("mug_shot")
	//c.Request.FormFile("file")
	if err != nil {
		sugar.Info(err.Error())
		c.JSON(handler.OK, handler.Response{
			ResultCode: handler.RequestFormatError1,
			Message:    handler.ResponseFlag[handler.RequestFormatError1],
			Data:       "",
			TimeStamp:  time.Now(),
		})
		return
	}
	gcp.InsertFileToGCS("mug-shot/", uploadFile, file)
	c.JSON(handler.OK, handler.Response{
		ResultCode: handler.OK,
		Message:    "上傳重工",
		Data:       request.MugShot.Filename,
		TimeStamp:  time.Now(),
	})
}

func GetMugShot(c *gin.Context) {
	tmp := gcp.GetFileForGCS("mug-shot/test.jpg")
	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": tmp,
	})
}
