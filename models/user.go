package models

import (
	"time"

	"gorm.io/gorm"
)

// User 对应数据库中的 users 表
type User struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`               // 主键，自增
	Username     string         `gorm:"type:VARCHAR(50);unique;not null" json:"username"` // 用户名，唯一且非空
	PasswordHash string         `gorm:"type:VARCHAR(255);not null" json:"-"`              // 密码哈希，非空（前端不返回）
	Role         string         `gorm:"type:VARCHAR(20);default:'user'" json:"role"`      // 角色，默认值为 'user'
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`                 // 创建时间，自动填充
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                 // 更新时间，自动更新
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`                // 软删除字段（仅当需要时显示）
}

// TableName 显式指定表名（与 SQL 中的表名保持一致）
func (User) TableName() string {
	return "users"
}
