package smtp

import (
	"context"
	"encoding/base64"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/pkg/logger"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()
var sugar = logger.Setup()
var redisClient = redis.Setup()

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

func sendEmailProcess(email string) string {
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
			"請點擊以下網址連結：" + os.Getenv("HOST_NAME") + "/api/v1/member/email-auth/" + emailToken + "\r\n" +
			"\r\n",
	)
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err := smtp.SendMail(os.Getenv("EMAIL_HOST")+":"+os.Getenv("EMAIL_PORT"), auth, from, to, message)
	if err != nil {
		sugar.Info(err.Error())
	}
	return emailToken
}

func insertEmailTokenToRedis(email string, emailToken string) {
	sugar.Info(email + "   " + emailToken)
	err := redisClient.Set(ctx, emailToken, email, 600*time.Second).Err()
	if err != nil {
		sugar.Info(err.Error())
	}
}

func SendEmailAuth(email string) {
	emailToken := sendEmailProcess(email)
	insertEmailTokenToRedis(email, emailToken)
}
