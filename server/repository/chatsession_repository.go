package repository

import (
	"gorm.io/gorm"
	"strix-server/common"
	"strix-server/persistence"
)

type ChatSessionRepositoryPostgres struct {
	DbContext *gorm.DB
}

func NewChatSessionRepository(context *gorm.DB) (u *ChatSessionRepositoryPostgres) {
	return &ChatSessionRepositoryPostgres{
		DbContext: context,
	}
}

func (u *ChatSessionRepositoryPostgres) FindById(ID string, target *persistence.ChatSession) error {
	uuID := common.GetUUIDFromString(ID)
	err := u.DbContext.Preload("Sender").Preload("Receiver").Where("id = ?", &uuID).First(target).Error
	return err
}

func (u *ChatSessionRepositoryPostgres) FindBySenderAndReciever(senderID string, recieverId string, target *persistence.ChatSession) error {
	uusenderID := common.GetUUIDFromString(senderID)
	uurecieverID := common.GetUUIDFromString(recieverId)
	err := u.DbContext.Preload("Sender").Preload("Receiver").Where("sender_id = ?", &uusenderID).Where("receiver_id = ?", &uurecieverID).First(target).Error
	return err
}

func (u *ChatSessionRepositoryPostgres) FindAllPending(userId string, target *[]persistence.ChatSession) error {
	userid := common.GetUUIDFromString(userId)
	err := u.DbContext.Preload("Sender").Preload("Receiver").Preload("Sender.PreKeys").Preload("Receiver.PreKeys").
		Where("receiver_id = ?", &userid).
		Where("is_initialized = ?", false).
		Find(target).
		Error
	return err
}

func (u *ChatSessionRepositoryPostgres) Save(target *persistence.ChatSession) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}
