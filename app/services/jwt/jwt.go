package jwt

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/junminhong/member-services-center/config/database"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()

var redisClient = database.InitRedis()

func getLocalSecretKey(fileName string) []byte {
	nowWorkDir, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	SECRETKEY, err := ioutil.ReadFile(nowWorkDir + "/" + fileName + ".pem")
	if err != nil {
		log.Println(err.Error())
	}
	return SECRETKEY
}

func GenerateAccessToken(memberID int) string {
	type MyCustomClaims struct {
		jwt.StandardClaims
	}
	now := time.Now()
	jwtID := strconv.FormatInt(now.Unix(), 10)
	claims := MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: now.Add(600 * time.Second).Unix(),
			Id:        jwtID,
			IssuedAt:  now.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(getLocalSecretKey("key"))
	if err != nil {
		log.Println(err.Error())
	}
	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		log.Println(err.Error())
	}
	//redisClient := database.InitRedis()
	err = redisClient.Set(ctx, accessToken, memberID, 600*time.Second).Err()
	if err != nil {
		log.Println(err.Error())
	}
	return accessToken
}

func VerifyAccessToken(accessToken string) bool {
	PUBKEY, _ := jwt.ParseRSAPublicKeyFromPEM(getLocalSecretKey("pubkey"))
	tokenParts := strings.Split(accessToken, ".")
	err := jwt.SigningMethodRS256.Verify(strings.Join(tokenParts[0:2], "."), tokenParts[2], PUBKEY)
	if err != nil {
		log.Println(err.Error())
	}
	type MyCustomClaims struct {
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(accessToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return PUBKEY, nil
	})
	if err != nil {
		log.Println(err.Error())
	}
	return token.Valid
}
