package controller

type Controllers struct {
	UserIer
	AdminIer
}

func New() *Controllers {
	return &Controllers{UserIer: &UserCtrl{}, AdminIer: &AdminCtrl{}}
}
func (c *Controllers) UserCtrl() UserIer {
	return c.UserIer
}
func (c *Controllers) AdminCtrl() AdminIer {
	return c.AdminIer
}
