package entity

type User struct {
	Model
	Mac    string `gorm:"column:f_mac;NOT NULL"`              // mac
	IP     string `gorm:"column:f_ip;NOT NULL"`               // ip地址
	Code   string `gorm:"column:f_code;NOT NULL"`             // 绑定code
	Use    int    `gorm:"column:f_use;default:0;NOT NULL"`    // 是否使用，0未使用，1使用,2未使用
	Online int    `gorm:"column:f_online;default:0;NOT NULL"` // 是否在线，1在线，2离线
}

func (m *User) TableName() string {
	return "t_user"
}
