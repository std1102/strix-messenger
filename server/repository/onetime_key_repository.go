package repository

import (
	"gorm.io/gorm"
	"strix-server/persistence"
)

type OneTimeKeyRepository struct {
	DbContext *gorm.DB
}

func NewOneTimeKeyRepository(context *gorm.DB) (u *OneTimeKeyRepository) {
	return &OneTimeKeyRepository{
		DbContext: context,
	}
}

func (u *OneTimeKeyRepository) Save(target *persistence.PreKeys) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}
