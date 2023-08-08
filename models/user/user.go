package user

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	Pk        uint       `gorm:"primary_key;auto_increment"` // 自增主键
	Id        string     `gorm:"column:id"`                  // 用户Id
	Username  string     `gorm:"column:username"`            // 用户名
	Password  string     `gorm:"column:password"`            // 密码
	CreatedAt time.Time  `gorm:"column:created_at"`          // 创建时间
	UpdatedAt time.Time  `gorm:"column:updated_at"`          // 更新时间
	DeletedAt *time.Time `gorm:"column:deleted_at"`          // 删除时间
}

func (t *User) Create(db *gorm.DB) error {
	return db.Create(t).Error
}
