package entity

import "time"

type Model struct {
	Id        int64      `gorm:"column:f_id;primaryKey" db:"f_id" json:"id" form:"f_id"`                                                          // 自增id
	CreatedAt *time.Time `gorm:"column:f_created_at;default:current_timestamp"  json:"createdAt" form:"f_created_at"`                             // 创建时间
	UpdatedAt *time.Time `gorm:"column:f_updated_at;default:current_timestamp on update current_timestamp"  json:"updatedAt" form:"f_updated_at"` // 更新时间
}
