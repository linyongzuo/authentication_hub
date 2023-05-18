package filter

import "gorm.io/gorm"

func WithMac(mac string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("f_mac = ? ", mac)
		return db
	}
}
func WithCode(code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("f_code = ? ", code)
		return db
	}
}
func WithMacCode(mac string, code string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("f_mac = ? and f_code = ?", mac, code)
		return db
	}
}
func WithIp(ip string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("f_ip", ip)
		return db
	}
}
func WithUsed() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = db.Where("f_use", 1)
		return db
	}
}
