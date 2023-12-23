package repository

import (
	"gorm.io/gorm"
	"strix-server/common"
	"strix-server/persistence"
)

type FileRepositoryPostgres struct {
	DbContext *gorm.DB
}

func NewFileRepository(context *gorm.DB) (u *FileRepositoryPostgres) {
	return &FileRepositoryPostgres{
		DbContext: context,
	}
}

func (u *FileRepositoryPostgres) FindById(ID string, target *persistence.UploadedFile) error {
	uuID := common.GetUUIDFromString(ID)
	err := u.DbContext.Preload("Owner").Where("id = ?", &uuID).First(target).Error
	return err
}

func (u *FileRepositoryPostgres) Save(target *persistence.UploadedFile) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}
