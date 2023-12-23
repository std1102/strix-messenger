package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	cmap "github.com/orcaman/concurrent-map/v2"
	"net/http"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
	"time"
)

type SocketSession struct {
	InitTime int64
	User     *persistence.User
}

type VOIPSession struct {
	InitTime     int64
	Sender       *persistence.User
	Reciever     *persistence.User
	CallerConn   *websocket.Conn
	RecieverConn *websocket.Conn
}

var CURRENT_USER_ACTIVE = cmap.New[*websocket.Conn]()
var SOCKET_SESSION_TOKEN = cmap.New[*SocketSession]()
var VOIP_SESSION_TOKEN = cmap.New[*VOIPSession]()

var myUpgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func cleanUpChatSocketSession() {
	for {
		currentTimeStamp := time.Now().UnixMilli()
		for k := range SOCKET_SESSION_TOKEN.Items() {
			tokenInfo, existed := SOCKET_SESSION_TOKEN.Get(k)
			if existed && tokenInfo != nil {
				if (tokenInfo.InitTime + 60000) < currentTimeStamp {
					SOCKET_SESSION_TOKEN.Remove(k)
				}
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func cleanUpVoipSocketSession() {
	for {
		currentTimeStamp := time.Now().UnixMilli()
		for k := range VOIP_SESSION_TOKEN.Items() {
			tokenInfo, existed := VOIP_SESSION_TOKEN.Get(k)
			if existed && tokenInfo != nil {
				if (tokenInfo.InitTime+15000) < currentTimeStamp && (tokenInfo.CallerConn == nil && tokenInfo.RecieverConn == nil) {
					VOIP_SESSION_TOKEN.Remove(k)
				}
			}
		}
		time.Sleep(15 * time.Second)
	}
}

// Communicate
func initSocketSession(context *gin.Context) {
	user := getLoggedInUser(context)
	rndBytes, _ := common.RandomBytes(32)
	randomToken := common.EncodeToString(rndBytes)
	SOCKET_SESSION_TOKEN.Set(randomToken, &SocketSession{
		InitTime: time.Now().UnixMilli(),
		User:     user,
	})
	context.JSON(200, gin.H{
		"authToken": randomToken,
	})
}

func initVoipSession(context *gin.Context) {
	recieverUserName := context.Query("userId")
	callType := context.Query("callType")
	ephemeralKey := context.Query("ephemeralKey")
	if recieverUserName == "" {
		handleError(context, 400, fmt.Errorf("Missing userId"))
		return
	}

	var recievedUser persistence.User
	userRepository := repository.NewUserRepository(persistence.DatabaseContext)
	err := userRepository.FindByUserName(recieverUserName, &recievedUser)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	currentUser := getLoggedInUser(context)

	// TODO Send notif
	// Check if user is active or not, if not reject
	// Else send for loop to send socket message to notif other user

	rndBytes, _ := common.RandomBytes(32)
	randomToken := common.EncodeToString(rndBytes)

	VOIP_SESSION_TOKEN.Set(randomToken, &VOIPSession{
		InitTime:     time.Now().UnixMilli(),
		Sender:       currentUser,
		Reciever:     &recievedUser,
		CallerConn:   nil,
		RecieverConn: nil,
	})

	sendCallingMessage(currentUser.Username, recievedUser.ID.String(), randomToken, callType, ephemeralKey)

	context.JSON(200, gin.H{
		"voipSession": randomToken,
	})
}

func connectVoipCall(context *gin.Context) {
	voipSession := context.Query("voipSession")
	tokenInfo, existed := VOIP_SESSION_TOKEN.Get(voipSession)
	if !existed {
		handleError(context, 401, fmt.Errorf("Unauthorized"))
		return
	}
	connType := context.Query("connType")
	if connType == "" {
		handleError(context, 400, fmt.Errorf("Missing connType"))
		return
	}

	wsConn, err := myUpgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if "FROM_CALLER" == connType {
			tokenInfo.CallerConn = nil
		} else if "FROM_RECIEVER" == connType {
			tokenInfo.RecieverConn = nil
		}
		if err != nil {
			if "FROM_CALLER" == connType {
				tokenInfo.CallerConn = nil
			} else if "FROM_RECIEVER" == connType {
				tokenInfo.RecieverConn = nil
			}
			system.Logger.Error("What the fuck can i do ", err)
		}
	}(wsConn)

	if "FROM_CALLER" == connType {
		tokenInfo.CallerConn = wsConn
	} else if "FROM_RECIEVER" == connType {
		tokenInfo.RecieverConn = wsConn
	}

	for {
		mt, msgData, err := wsConn.ReadMessage()
		if err != nil {
			system.Logger.Error(err)
			continue
		}
		tokenInfo, existed = VOIP_SESSION_TOKEN.Get(voipSession)
		if !existed {
			handleError(context, 400, fmt.Errorf("Session ended"))
			return
		}
		if "FROM_CALLER" == connType {
			err := tokenInfo.RecieverConn.WriteMessage(mt, msgData)
			if err != nil {
				system.Logger.Infof("RECIEVER ", err)
				continue
			}
		} else if "FROM_RECIEVER" == connType {
			err := tokenInfo.CallerConn.WriteMessage(mt, msgData)
			if err != nil {
				system.Logger.Infof("RECIEVER ", err)
				continue
			}
		}
	}
}

func webSocket(context *gin.Context) {
	authToken := context.Query("authToken")
	tokenInfo, existed := SOCKET_SESSION_TOKEN.Get(authToken)

	if !existed {
		handleError(context, 401, fmt.Errorf("Unauthorized"))
		return
	} else {
		SOCKET_SESSION_TOKEN.Remove(authToken)
	}

	currentUser := tokenInfo.User

	conn, err := myUpgrader.Upgrade(context.Writer, context.Request, nil)

	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}

	CURRENT_USER_ACTIVE.Set(currentUser.ID.String(), conn)

	chatSessionRepository := repository.NewChatSessionRepository(persistence.DatabaseContext)
	pendingMessageRepository := repository.NewPendingMessageRepository(persistence.DatabaseContext)

	cachedConversation := make(map[string]*persistence.ChatSession)

	defer func(conn *websocket.Conn) {
		CURRENT_USER_ACTIVE.Set(currentUser.ID.String(), nil)
		err := conn.Close()
		if err != nil {
			system.Logger.Error("What the fuck can i do ", err)
		}
	}(conn)

	for {
		// read msgData
		mt, msgData, err := conn.ReadMessage()
		if err != nil {
			system.Logger.Error(err)
			continue
		}
		var msgDto MessageDto
		err = json.Unmarshal(msgData, &msgDto)
		if err != nil {
			// TODO handle error here
			fmt.Println(err.Error())
			system.Logger.Error(err)
			continue
		}

		msgDto.SenderUsername = currentUser.Username

		targetChatSession := cachedConversation[msgDto.ChatSessionId]

		msgData, err = json.Marshal(&msgDto)
		if err != nil {
			// TODO handle error here
			system.Logger.Error(err)
			continue
		}

		if msgDto.Type == CHAT_ACCEPT || msgDto.Type == CHAT_CLOSE {
			var recievedUser persistence.User
			userRepository := repository.NewUserRepository(persistence.DatabaseContext)
			err := userRepository.FindByUserName(*msgDto.PlainMessage, &recievedUser)
			if err != nil {
				system.Logger.Errorf(err.Error())
				continue
			}
			otherConn, _ := CURRENT_USER_ACTIVE.Get(recievedUser.ID.String())
			err = otherConn.WriteMessage(mt, msgData)
			if err != nil {
				system.Logger.Errorf(err.Error())
				continue
			}
			continue
		}

		if targetChatSession != nil {
			fromSender := targetChatSession.SenderId.String() == currentUser.ID.String()
			result := sendMessage(mt, &msgDto, fromSender, targetChatSession, pendingMessageRepository)
			if result == 0 {

			} else if result == 2 {
				// If user not active
			}
			continue
		} else {
			var chatSessionInDb persistence.ChatSession
			err := chatSessionRepository.FindById(msgDto.ChatSessionId, &chatSessionInDb)
			if err != nil {
				system.Logger.Error(err)
				continue
			}
			fromSender := chatSessionInDb.SenderId.String() == currentUser.ID.String()
			cachedConversation[msgDto.ChatSessionId] = &chatSessionInDb
			result := sendMessage(mt, &msgDto, fromSender, &chatSessionInDb, pendingMessageRepository)
			if result == 0 {
				// TODO Handle error
				continue
			} else if result == 2 {
				// if user not active
			}
			continue
		}
	}
}

func sendMessage(mt int, msgDto *MessageDto, fromSender bool, chatSession *persistence.ChatSession, pendingMessageRepository *repository.PendingMessageRepositoryPostgres) int {
	var targetUser *persistence.User
	if fromSender {
		targetUser = chatSession.Receiver
	} else {
		targetUser = chatSession.Sender
	}
	otherConn, existed := CURRENT_USER_ACTIVE.Get(targetUser.ID.String())
	msgData, _ := json.Marshal(msgDto)
	if existed && otherConn != nil {
		// if user is online
		err := otherConn.WriteMessage(mt, msgData)
		if err != nil {
			system.Logger.Error(err)
			err := savePendingMessage(msgDto, fromSender, chatSession, pendingMessageRepository)
			if err != nil {
				system.Logger.Error(err)
				return 0
			}
			return 0
		}
		return 1
	} else {
		// if user is offline - save message to pending message
		err := savePendingMessage(msgDto, fromSender, chatSession, pendingMessageRepository)
		if err != nil {
			return 0
		}
		return 2
	}
}

func savePendingMessage(msg *MessageDto, fromSender bool, chatSession *persistence.ChatSession, pendingMessageRepository *repository.PendingMessageRepositoryPostgres) error {
	pendingId, _ := uuid.NewUUID()
	var owner *persistence.User
	var sender *persistence.User
	if fromSender {
		owner = chatSession.Receiver
		sender = chatSession.Sender
	} else {
		owner = chatSession.Sender
		sender = chatSession.Receiver
	}

	pendingMessage := persistence.PendingMessage{
		ID:             pendingId,
		Type:           msg.Type,
		Index:          msg.Index,
		OwnerId:        owner.ID,
		SenderId:       sender.ID,
		SenderUsername: msg.SenderUsername,
		ChatSessionId:  chatSession.ID,
		CipherMessage:  msg.CipherMessage,
		PlainMessage:   msg.PlainMessage,
		IsBinary:       msg.IsBinary,
		IsRead:         false,
		Owner:          owner,
		Sender:         sender,
		ChatSession:    chatSession,
		CreatedAt:      time.Now(),
	}
	err := pendingMessageRepository.Insert(&pendingMessage)
	if err != nil {
		system.Logger.Error(err)
		return err
	}
	return nil
}

func sendCallingMessage(senderUsername, recvUsername string, voipToken, callType, ephemeralKey string) {
	msgDto := MessageDto{
		Type:           callType,
		PlainMessage:   &ephemeralKey,
		CipherMessage:  voipToken,
		SenderUsername: senderUsername,
	}
	msgBin, err := json.Marshal(&msgDto)
	if err != nil {
		system.Logger.Error(err)
		return
	}
	otherConn, existed := CURRENT_USER_ACTIVE.Get(recvUsername)
	if !existed || otherConn == nil {
		return
	}
	err = otherConn.WriteMessage(1, msgBin)
	if err != nil {
		return
	}
	/*for i := 0; i < 30; i++ {
		err = otherConn.WriteMessage(1, msgBin)
		if err != nil {
			continue
		}
		time.Sleep(1 * time.Second)
	}*/
}
