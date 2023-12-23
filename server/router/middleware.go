package router

import (
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
	"time"
)

// Middlewares

func CORSMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Allow-Origin")
	c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")

	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(204)
		return
	}

	c.Next()
}

func contentTypeMiddleware(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Next()
}

func loggingMiddleWare(context *gin.Context) {
	requestId := requestid.Get(context)
	context.Header(REQUEST_ID, requestId)
	system.Logger.Infof(
		"app-node: %s | %s:%s | from: %s | request-id: %s",
		system.SystemConfig.App.Node,
		context.Request.Method,
		context.FullPath(),
		context.ClientIP(),
		requestId)
	context.Next()
}

func errorHandlerMiddleWare(context *gin.Context) {
	requestId := requestid.Get(context)
	context.Next()

	t := time.Now()
	timeString := common.FormatTime(&t)

	for status, err := range context.Errors {
		system.Logger.Errorf("RequestId %s | Error: %s", requestId, err.Error())
		context.JSON(status, gin.H{
			"requestId": requestId,
			"message":   err.Error(),
			"time":      timeString,
		})
	}

}

func authenticationMiddleWare(context *gin.Context) {
	path := context.Request.URL.Path
	if path == "/api/v1/auth/register" || path == "/api/v1/auth/login" || path == "/ws" || path == "/voip" || path == "/api/v1/file/get" {
		context.Next()
		return
	}
	token := context.GetHeader("Authorization")
	if token == "" {
		handleError(context, 401, fmt.Errorf("Unauthroized"))
		context.Next()
		return
	}
	token = token[7:]
	var claimMap jwt.MapClaims
	_, err := jwt.ParseWithClaims(token, &claimMap, func(token *jwt.Token) (any, error) {
		return common.DecodeToByte(system.SystemConfig.JwtKey), nil
	})
	if err != nil {
		handleError(context, 401, fmt.Errorf("Unauthroized"))
		context.Next()
		return
	}
	userId := claimMap["userId"].(string)
	expiredAt := claimMap["exp"].(float64)
	currentTime := time.Now()
	if float64(currentTime.UnixMilli()) > expiredAt {
		handleError(context, 401, fmt.Errorf("Unauthroized"))
		context.Next()
		return
	}
	var userRepository = repository.NewUserRepository(persistence.DatabaseContext)
	var user persistence.User
	err = userRepository.FindById(userId, &user)
	if err != nil {
		handleError(context, 401, fmt.Errorf(err.Error()))
		context.Next()
		return
	}
	context.Set(USER, &user)
	system.Logger.Infof("User: %s logged in at: %s", user.Username, common.FormatTime(&currentTime))
	context.Next()
}
