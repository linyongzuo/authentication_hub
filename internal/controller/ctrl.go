package controller

import (
	"github.com/authentication_hub/global"
	"github.com/authentication_hub/internal/repository/persistence/impl"
)

type Controllers struct {
	UserIer
	AdminIer
}

func New() *Controllers {
	return &Controllers{UserIer: &UserCtrl{db: global.DbPool, userIer: impl.NewUser()}, AdminIer: &AdminCtrl{db: global.DbPool, adminIer: impl.NewAdmin(), userIer: impl.NewUser()}}
}
func (c *Controllers) UserCtrl() UserIer {
	return c.UserIer
}
func (c *Controllers) AdminCtrl() AdminIer {
	return c.AdminIer
}
