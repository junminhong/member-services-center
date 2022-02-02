package smtp

import (
	"context"
	"encoding/base64"
	"github.com/joho/godotenv"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"os"
	"time"
)

var ctx = context.Background()
var sugar = logger.Setup()
var redisClient = redis.Setup()

const (
	// minutes為單位
	emailTokenTimeLimit = 60
)

func init() {
	err := godotenv.Load()
	if err != nil {
		sugar.Info(err.Error())
	}
}

func encodeEmailToken(email string) string {
	pwd := []byte(email)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		sugar.Info(err.Error())
	}
	return base64.StdEncoding.EncodeToString([]byte(hash))
}

func decodeEmailToken(token string) {
	base64.StdEncoding.DecodeString(token)
}

func sendEmailProcess(email string) bool {
	from := os.Getenv("EMAIL_ACCOUNT")
	password := os.Getenv("EMAIL_PASSWORD")
	to := []string{
		email,
	}
	emailToken := encodeEmailToken(email)
	message := []byte(
		"Subject: 帳戶認證信件\r\n" +
			"From: membercentersmtp@gmail.com\r\n" +
			`Content-Type: text/plain; boundary="qwertyuio"` + "\r\n" +
			"\r\n" +
			"請點擊以下網址連結：" + os.Getenv("HOST_NAME") + os.Getenv("EMAIL_AUTH_PATH") + "?email_token=" + emailToken + "\r\n" +
			"\r\n",
	)
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"), auth, from, to, message)
	if err != nil {
		sugar.Info(err.Error())
		return false
	}
	sugar.Info("已將驗證信件寄送給：" + email)
	return insertEmailTokenToRedis(email, emailToken)
}

func insertEmailTokenToRedis(email string, emailToken string) bool {
	err := redisClient.Set(ctx, emailToken, email, emailTokenTimeLimit*time.Minute).Err()
	if err != nil {
		sugar.Info(err.Error())
		return false
	}
	sugar.Info("已將" + emailToken + "存進redis，有效期為" + (emailTokenTimeLimit * time.Minute).String())
	return true
}

func SendEmailAuth(email string) bool {
	sugar.Info("正在寄送帳戶驗證信件")
	return sendEmailProcess(email)
}
