package router

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/system"
	"time"
)

func handleError(context *gin.Context, statusCode int, err error) {
	system.Logger.Error(err.Error())
	e := context.AbortWithError(statusCode, err)
	if e != nil {
		system.Logger.Error(err.Error())
	}
	context.Next()
}

func generateToken(initTime *time.Time, expiredTime uint64, userId string) (string, error) {
	iatTime := initTime.UnixMilli()
	timeToTLive := int64(expiredTime)
	expTime := iatTime + timeToTLive

	jwtKey := common.DecodeToByte(system.SystemConfig.JwtKey)
	jwtProp := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userId,
		"iat":    iatTime,
		"ttl":    timeToTLive,
		"exp":    expTime,
	})
	return jwtProp.SignedString(jwtKey)
}

func getLoggedInUser(context *gin.Context) *persistence.User {
	u, isExist := context.Get(USER)
	if !isExist {
		return nil
	}
	return u.(*persistence.User)
}
