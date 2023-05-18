package impl

import (
	"github.com/authentication_hub/internal/repository/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AdminIerImpl struct {
}

func NewAdmin() *AdminIerImpl {
	return &AdminIerImpl{}
}
func (a *AdminIerImpl) BatchSave(db *gorm.DB, admins []*entity.Admin) error {
	err := db.CreateInBatches(admins, 1000).Error
	if err != nil {
		logrus.Errorf("批量创建账号失败:%s", err.Error())
		return err
	}
	return nil
}
func (a *AdminIerImpl) SearchAdmin(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (admins []*entity.Admin, err error) {
	db = db.Model(&entity.Admin{}).Scopes(scopes...)
	err = db.Find(&admins).Error
	return
}

func (a *AdminIerImpl) Get(tx *gorm.DB, in *entity.Admin, scopes ...func(*gorm.DB) *gorm.DB) (out *entity.Admin, err error) {
	err = tx.Model(&entity.Admin{}).Scopes(scopes...).Where(&in).First(&out).Error
	if err != nil {
		return
	}
	return
}

func (a *AdminIerImpl) Count(tx *gorm.DB, in *entity.Admin, scopes ...func(*gorm.DB) *gorm.DB) (count int64, err error) {
	err = tx.Model(&entity.Admin{}).Scopes(scopes...).Where(&in).Count(&count).Error
	if err != nil {
		return
	}
	return
}
