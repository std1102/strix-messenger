package repository

import (
	"gorm.io/gorm"
	"strix-server/common"
)

type DataSource interface {
	FindById(ID any) any
	FindAll() any
	FindAllPaging(pageable *common.Pageable) any
	Insert(obj any)
	Delete(ID any)
}

func GetPageable(page *common.Pageable, context *gorm.DB, bindValues any) error {
	getPageable(page)
	return context.
		Preload("OneTimePreKeys").
		Limit(int(page.PageSize)).
		Offset(int(page.PageSize * page.PageNumber)).
		Find(bindValues).Error
}

func getPageable(page *common.Pageable) {
	if page == nil {
		page = &common.Pageable{
			PageNumber: 1,
			PageSize:   25,
		}
	} else {
		if page.PageSize <= 0 || page.PageSize > 25 {
			page.PageSize = 25
		}
		if page.PageNumber <= 0 {
			page.PageNumber = 1
		}
	}
}
