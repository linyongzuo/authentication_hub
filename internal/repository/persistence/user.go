package persistence

import (
	"github.com/authentication_hub/internal/repository/entity"
	"gorm.io/gorm"
)

type UserIer interface {
	BatchSave(db *gorm.DB, users []*entity.User) error
	SearchUser(db *gorm.DB, scopes ...func(*gorm.DB) *gorm.DB) (users []*entity.User, err error)
	Update(db *gorm.DB, in *entity.User, scopes ...func(db *gorm.DB) *gorm.DB) (out *entity.User, err error)
	Get(tx *gorm.DB, in *entity.User, scopes ...func(*gorm.DB) *gorm.DB) (out *entity.User, err error)
}
