package models

import (
	"time"

	"gorm.io/gorm"
)

// Paper 对应数据库中的 papers 表
type Paper struct {
	ID          int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Title       string         `gorm:"type:VARCHAR(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	TotalScore  int            `gorm:"default:100" json:"total_score"`
	CreatorID   int64          `gorm:"not null" json:"creator_id"`
	Creator     User           `gorm:"foreignKey:CreatorID" json:"creator,omitempty"` // 关联创建者
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 显式指定表名
func (Paper) TableName() string {
	return "papers"
}
