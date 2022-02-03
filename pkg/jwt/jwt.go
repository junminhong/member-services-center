package jwt

import (
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/junminhong/member-services-center/db/redis"
	"github.com/junminhong/member-services-center/pkg/logger"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var ctx = context.Background()
var sugar = logger.Setup()
var redisClient = redis.Setup()

const (
	// atomic token有效期限分鐘為單位
	atomicTokenLimitTime = 60
	// RefreshAtomicTokenExpired refresh atomic token有效期限小時為單位
	RefreshAtomicTokenExpired = 24
)

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

func GenerateRefreshAtomicToken(memberUUID string) string {
	now := time.Now()
	jwtID := strconv.FormatInt(now.Unix(), 10)
	claims := &jwt.StandardClaims{
		ExpiresAt: now.Add(RefreshAtomicTokenExpired * time.Hour).Unix(),
		Id:        jwtID,
		IssuedAt:  now.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(getLocalSecretKey("key"))
	if err != nil {
		sugar.Info(err.Error())
	}
	refreshAtomicToken, err := token.SignedString(privateKey)
	if err != nil {
		sugar.Info(err.Error())
	}
	err = redisClient.Set(ctx, refreshAtomicToken, memberUUID, RefreshAtomicTokenExpired*time.Hour).Err()
	if err != nil {
		sugar.Info(err.Error())
	}
	return refreshAtomicToken
}

func GenerateAccessToken(memberUUID string) string {
	type MyCustomClaims struct {
		jwt.StandardClaims
	}
	now := time.Now()
	jwtID := strconv.FormatInt(now.Unix(), 10)
	claims := MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: now.Add(atomicTokenLimitTime * time.Minute).Unix(),
			Id:        jwtID,
			IssuedAt:  now.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(getLocalSecretKey("key"))
	if err != nil {
		sugar.Info(err.Error())
	}
	atomicToken, err := token.SignedString(privateKey)
	if err != nil {
		sugar.Info(err.Error())
	}
	//redisClient := database.InitRedis()
	err = redisClient.Set(ctx, atomicToken, memberUUID, atomicTokenLimitTime*time.Minute).Err()
	if err != nil {
		sugar.Info(err.Error())
	}
	return atomicToken
}

func VerifyAtomicToken(atomicToken string) bool {
	PUBKEY, _ := jwt.ParseRSAPublicKeyFromPEM(getLocalSecretKey("pubkey"))
	tokenParts := strings.Split(atomicToken, ".")
	err := jwt.SigningMethodRS256.Verify(strings.Join(tokenParts[0:2], "."), tokenParts[2], PUBKEY)
	if err != nil {
		sugar.Info(err.Error())
	}
	type MyCustomClaims struct {
		jwt.StandardClaims
	}
	token, err := jwt.ParseWithClaims(atomicToken, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return PUBKEY, nil
	})
	if err != nil {
		sugar.Info(err.Error())
	}
	return token.Valid
}
