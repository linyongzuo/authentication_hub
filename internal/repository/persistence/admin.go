package persistence

import (
	"github.com/authentication_hub/internal/repository/entity"
	"gorm.io/gorm"
)

type AdminIer interface {
	BatchSave(db *gorm.DB, admins []*entity.Admin) error
	SearchAdmin(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (admins []*entity.Admin, err error)
	Get(tx *gorm.DB, in *entity.Admin, scopes ...func(*gorm.DB) *gorm.DB) (out *entity.Admin, err error)
}
