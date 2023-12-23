package router

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strix-server/system"
	"time"
)

const REQUEST_ID = "X-Request-ID"
const BCRYPT_COST = 12
const USER = "user"

var router *gin.Engine

var upgrader = websocket.Upgrader{}

func Init() {
	go cleanUpChatSocketSession()
	router = gin.New()
	// Middleware
	router.Use(
		requestid.New(
			requestid.WithGenerator(func() string {
				requestId, _ := uuid.NewUUID()
				return requestId.String()
			}),
			requestid.WithCustomHeaderStrKey("X-Request-ID"),
		),
	)
	router.Use(errorHandlerMiddleWare)
	router.Use(loggingMiddleWare)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     system.SystemConfig.Server.AllowOrigins,
		AllowMethods:     system.SystemConfig.Server.AllowMethods,
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(contentTypeMiddleware)
	router.Use(authenticationMiddleWare)

	// Swagger
	router.GET("/api-doc", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routing
	testGroup := router.Group("/api/v1/hello")
	testGroup.GET("/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "hello server",
		})
	})

	// Auth API
	authenticaionGroup := router.Group("/api/v1/auth")
	authenticaionGroup.POST("/register", register)
	authenticaionGroup.POST("/login", login)

	// User API
	userGroup := router.Group("/api/v1/user")
	userGroup.POST("/uploadKey", uploadKey)
	userGroup.GET("/:userName/externalKey", getExternalKeyBundle)
	userGroup.GET("", getUserInfo)
	userGroup.GET("/search", searchUser)
	userGroup.POST("/userInfos", getUserInfos)

	// Chat Session API
	chatSessionGroup := router.Group("/api/v1/chatSession")
	chatSessionGroup.POST("/init", initChatSession)
	chatSessionGroup.GET("", retrieveChatSession)
	chatSessionGroup.PUT("/complete", completeInitChatSession)

	// Message
	messageGroup := router.Group("/api/v1/message")
	messageGroup.GET("", retrievePendingMessage)

	// File
	fileGroup := router.Group("/api/v1/file")
	fileGroup.POST("", uploadFile)
	fileGroup.POST("/avatar", uploadAvatar)
	fileGroup.GET("/get", getFile)

	// Communication
	router.GET("/api/v1/ws/init", initSocketSession)
	router.PUT("/api/v1/voip/init", initVoipSession)

	// Ws
	router.GET("/ws", webSocket)
	router.GET("/voip", connectVoipCall)

	err := router.Run(fmt.Sprintf(":%s", system.SystemConfig.Server.Port))
	if err != nil {
		logrus.Fatal("Cannot start server", err)
	}
	go system.Logger.Infof("Started server and listening at: %s:%s", system.SystemConfig.Server.Address, system.SystemConfig.Server.Port)
}
