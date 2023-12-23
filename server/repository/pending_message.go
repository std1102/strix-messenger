package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"strix-server/common"
	"strix-server/persistence"
)

type PendingMessageRepositoryPostgres struct {
	DbContext *gorm.DB
}

func NewPendingMessageRepository(context *gorm.DB) (u *PendingMessageRepositoryPostgres) {
	return &PendingMessageRepositoryPostgres{
		DbContext: context,
	}
}

func (u *PendingMessageRepositoryPostgres) FindByUserNameAndChatSession(userId string, chatSessionId string, target *[]persistence.PendingMessage) error {
	userid := common.GetUUIDFromString(userId)
	chatsessionid := common.GetUUIDFromString(chatSessionId)
	err := u.DbContext.Where("owner_id", &userid).Where("chat_session_id", &chatsessionid).Find(target).Error
	return err
}

func (u *PendingMessageRepositoryPostgres) Insert(target *persistence.PendingMessage) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Create(target).Error
		return err
	})
}

func (u *PendingMessageRepositoryPostgres) Save(target *persistence.PendingMessage) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}

func (u *PendingMessageRepositoryPostgres) SaveAll(target *[]persistence.PendingMessage) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}

func (u *PendingMessageRepositoryPostgres) DeleteAll(target *[]persistence.PendingMessage, ids []uuid.UUID) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Delete(target, ids).Error
		return err
	})
}
