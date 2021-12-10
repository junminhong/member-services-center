package smtp

import (
	"context"
	"encoding/base64"
	"log"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/junminhong/member-services-center/config/database"
	"golang.org/x/crypto/bcrypt"
)

var ctx = context.Background()

func encodeEmailToken(email string) string {
	pwd := []byte(email)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err.Error())
	}
	return base64.StdEncoding.EncodeToString([]byte(hash))
}

func decodeEmailToken(token string) {
	base64.StdEncoding.DecodeString(token)
}

func sendEmailProcess(email string) string {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
	}
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
			"請點擊以下網址連結：http://127.0.0.1:8080/auth/" + emailToken + "\r\n" +
			"\r\n",
	)
	auth := smtp.PlainAuth("", from, password, "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, from, to, message)
	if err != nil {
		log.Println(err.Error())
	}
	return emailToken
}

func insertEmailTokenToRedis(email string, emailToken string) {
	redisClient := database.InitRedis()
	log.Println(emailToken, email)
	err := redisClient.Set(ctx, emailToken, email, 600*time.Second).Err()
	if err != nil {
		log.Println(err.Error())
	}
}

func SendEmailAuth(email string) {
	emailToken := sendEmailProcess(email)
	insertEmailTokenToRedis(email, emailToken)
}
