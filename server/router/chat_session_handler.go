package router

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"strix-server/common"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
	"time"
)

// Chat Session
func initChatSession(context *gin.Context) {
	var chatSessionDto ChatSessionDto
	err := context.BindJSON(&chatSessionDto)
	if err != nil {
		handleError(context, 400, fmt.Errorf(err.Error()))
		return
	}

	currentUser := getLoggedInUser(context)
	var otherUser persistence.User

	userRepository := repository.NewUserRepository(persistence.DatabaseContext)

	err = userRepository.FindByUserName(chatSessionDto.ReceiverUserName, &otherUser)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	chatSessionRepository := repository.NewChatSessionRepository(persistence.DatabaseContext)

	var checkChatSession persistence.ChatSession
	err = chatSessionRepository.FindBySenderAndReciever(currentUser.ID.String(), otherUser.ID.String(), &checkChatSession)
	if err == nil {
		handleError(context, 400, fmt.Errorf("Chat session existed"))
		return
	}

	newChatSession := persistence.ChatSession{
		ID:            common.GetUUIDFromString(chatSessionDto.ChatSessionId),
		SenderId:      currentUser.ID,
		ReceiverId:    otherUser.ID,
		IsInitialized: false,
		EphemeralKey:  chatSessionDto.EphemeralKey,
		DeletedAt:     nil,
		CreatedAt:     time.Now(),
		Sender:        currentUser,
		Receiver:      &otherUser,
	}

	err = chatSessionRepository.Save(&newChatSession)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}

	otherConn, existed := CURRENT_USER_ACTIVE.Get(otherUser.ID.String())
	if existed && otherConn != nil {
		helloMessage := fmt.Sprintf("User %s want to chat with you", currentUser.Username)
		var lastedOneTimeKey persistence.OneTimeKey
		for _, element := range currentUser.OneTimePreKeys {
			lastedOneTimeKey = *element
		}
		msg := MessageDto{
			Type:           CHAT_NEW,
			SenderUsername: currentUser.Username,
			PlainMessage:   &helloMessage,
			ChatSessionId:  newChatSession.ID.String(),
			Index:          0,
			CipherMessage:  "",
			IsBinary:       false,
			AdditionalData: &ChatSessionDto{
				ChatSessionId:    newChatSession.ID.String(),
				EphemeralKey:     newChatSession.EphemeralKey,
				ReceiverUserName: currentUser.Username,
				SenderUserName:   otherUser.Username,
				SenderKeyBundle: ExternalKeyBundleDto{
					IdentityKey:   currentUser.IdentityKey,
					PreKey:        currentUser.PreKey,
					PreKeySig:     currentUser.PreKeySignature,
					OneTimeKeyId:  lastedOneTimeKey.ID.String(),
					OneTimeKey:    lastedOneTimeKey.Key,
					OneTimeKeySig: lastedOneTimeKey.KeySignature,
				},
			},
		}

		binMsg, err := json.Marshal(&msg)
		if err != nil {
			system.Logger.Error(err)
		}

		// if user is online
		err = otherConn.WriteMessage(websocket.TextMessage, binMsg)
		if err != nil {
			system.Logger.Error(err)
		}

		system.Logger.Infof("Send new chat notif")
	}

	context.JSON(200, gin.H{
		"message": "ok",
	})
}

func retrieveChatSession(context *gin.Context) {
	currentUser := getLoggedInUser(context)
	chatSessionRepository := repository.NewChatSessionRepository(persistence.DatabaseContext)
	var chatSessionList []persistence.ChatSession
	err := chatSessionRepository.FindAllPending(currentUser.ID.String(), &chatSessionList)
	if err != nil {
		if err.Error() == "empty slice found" {
			context.JSON(200, make([]ChatSessionDto, 0))
			return
		}
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	var result []ChatSessionDto
	for i := range chatSessionList {
		currentChatSession := chatSessionList[i]
		sender := currentChatSession.Sender
		reciever := currentChatSession.Receiver
		var lastedOneTimeKey persistence.OneTimeKey
		for _, element := range sender.OneTimePreKeys {
			lastedOneTimeKey = *element
		}
		result = append(result, ChatSessionDto{
			ChatSessionId:    currentChatSession.ID.String(),
			EphemeralKey:     currentChatSession.EphemeralKey,
			ReceiverUserName: reciever.Username,
			SenderUserName:   sender.Username,
			SenderKeyBundle: ExternalKeyBundleDto{
				IdentityKey:   sender.IdentityKey,
				PreKey:        sender.PreKey,
				PreKeySig:     sender.PreKeySignature,
				OneTimeKeyId:  lastedOneTimeKey.ID.String(),
				OneTimeKey:    lastedOneTimeKey.Key,
				OneTimeKeySig: lastedOneTimeKey.KeySignature,
			},
		})
	}
	if result == nil {
		result = make([]ChatSessionDto, 0)
	}
	context.JSON(200, result)
}

func completeInitChatSession(context *gin.Context) {
	chatSessionId := context.Query("chatSessionId")
	if chatSessionId == "" {
		handleError(context, 400, fmt.Errorf("Missing chat session id"))
		return
	}
	var chatSession persistence.ChatSession
	chatSessionRepository := repository.NewChatSessionRepository(persistence.DatabaseContext)
	err := chatSessionRepository.FindById(chatSessionId, &chatSession)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	if chatSession.IsInitialized {
		handleError(context, 400, fmt.Errorf("Chat session initialized"))
		return
	}
	chatSession.IsInitialized = true
	err = chatSessionRepository.Save(&chatSession)
	if err != nil {
		handleError(context, 500, fmt.Errorf(err.Error()))
		return
	}
	context.JSON(200, gin.H{
		"ephemeralKey": chatSession.EphemeralKey,
	})
}
