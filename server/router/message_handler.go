package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strix-server/persistence"
	"strix-server/repository"
	"strix-server/system"
)

// Chat
func retrievePendingMessage(context *gin.Context) {
	chatSessionId := context.Query("chatSessionId")
	pendingMessageRepo := repository.NewPendingMessageRepository(persistence.DatabaseContext)
	var pendingMessages []persistence.PendingMessage
	currentUser := getLoggedInUser(context)
	err := pendingMessageRepo.FindByUserNameAndChatSession(currentUser.ID.String(), chatSessionId, &pendingMessages)
	if err != nil {
		system.Logger.Error(err)
		handleError(context, 500, fmt.Errorf("Internal error"))
		return
	}

	var result []MessageDto
	var deletedIds []uuid.UUID
	for i := range pendingMessages {
		currentMsg := pendingMessages[i]
		result = append(result, MessageDto{
			Type:           currentMsg.Type,
			SenderUsername: currentMsg.SenderUsername,
			PlainMessage:   currentMsg.PlainMessage,
			ChatSessionId:  currentMsg.ChatSessionId.String(),
			FilePath:       currentMsg.FilePath,
			Index:          currentMsg.Index,
			CipherMessage:  currentMsg.CipherMessage,
			IsBinary:       currentMsg.IsBinary,
		})
		currentMsg.IsRead = true
		deletedIds = append(deletedIds, currentMsg.ID)
	}
	if result != nil {
		err = pendingMessageRepo.DeleteAll(&pendingMessages, deletedIds)
	}
	if err != nil {
		if err.Error() == "empty slice found" {
			context.JSON(200, make([]MessageDto, 0))
			return
		}
		handleError(context, 500, fmt.Errorf("Internal error"))
		return
	}
	context.JSON(200, result)
}
