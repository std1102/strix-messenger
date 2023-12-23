package repository

import (
	"gorm.io/gorm"
	"strix-server/common"
	"strix-server/persistence"
)

type UserRepositoryPostgres struct {
	DbContext *gorm.DB
}

func NewUserRepository(context *gorm.DB) (u *UserRepositoryPostgres) {
	return &UserRepositoryPostgres{
		DbContext: context,
	}
}

func (u *UserRepositoryPostgres) FindById(ID string, target *persistence.User) error {
	uuID := common.GetUUIDFromString(ID)
	err := u.DbContext.Preload("OneTimePreKeys").First(target, "id", &uuID).Error
	return err
}

func (u *UserRepositoryPostgres) FindByUserName(userName string, target *persistence.User) error {
	err := u.DbContext.Preload("OneTimePreKeys").First(target, "username", userName).Error
	return err
}

func (u *UserRepositoryPostgres) FindAllByUserNames(userNames []string, target *[]persistence.User) error {
	err := u.DbContext.Preload("OneTimePreKeys").Where("username IN ?", userNames).Find(target).Error
	return err
}

func (u *UserRepositoryPostgres) SearchUserName(userName string, target *[]persistence.User) error {
	err := u.DbContext.Preload("OneTimePreKeys").Where("username like ?", "%"+userName+"%").Find(target).Error
	return err
}

func (u *UserRepositoryPostgres) FindAll(target *[]persistence.User) error {
	err := u.DbContext.Preload("OneTimePreKeys").Find(target).Error
	return err
}

func (u *UserRepositoryPostgres) FindAllPaging(target *[]persistence.User, pageable *common.Pageable) error {
	err := GetPageable(pageable, u.DbContext, target)
	return err
}

func (u *UserRepositoryPostgres) Save(target *persistence.User) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Save(target).Error
		return err
	})
}

/*func (u *UserRepositoryPostgres) Delete(ID string) error {
	return u.DbContext.Transaction(func(context *gorm.DB) error {
		err := u.DbContext.Delete(common.GetUUIDFromString(ID)).Error
		return err
	})
}*/
