package impl

import (
	"github.com/authentication_hub/internal/repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserIerImpl struct {
}

func NewUser() *UserIerImpl {
	return &UserIerImpl{}
}
func (u *UserIerImpl) BatchSave(db *gorm.DB, users []*entity.User) error {
	err := db.CreateInBatches(users, 1000).Error
	if err != nil {
		logrus.Errorf("批量创建账号失败:%s", err.Error())
		return err
	}
	return nil
}
func (u *UserIerImpl) SearchUser(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (users []*entity.User, err error) {
	db = db.Model(&entity.User{}).Scopes(scopes...)
	err = db.Find(&users).Error
	return
}
func (u *UserIerImpl) Update(db *gorm.DB, in *entity.User, scopes ...func(db *gorm.DB) *gorm.DB) (out *entity.User, err error) {
	err = db.Model(&entity.User{}).Scopes(scopes...).Updates(&in).Error
	return
}

func (u *UserIerImpl) Get(tx *gorm.DB, in *entity.User, scopes ...func(*gorm.DB) *gorm.DB) (out *entity.User, err error) {
	err = tx.Model(&entity.User{}).Scopes(scopes...).Where(&in).First(&out).Error
	return
}
