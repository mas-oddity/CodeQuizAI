package models

import (
	"time"

	"gorm.io/gorm"
)

// Question 对应数据库中的 questions 表
type Question struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string         `gorm:"type:text;not null" json:"title"`
	QuestionType string         `gorm:"type:VARCHAR(20);not null" json:"question_type"` // 'single' 或 'multiple'
	Options      string         `gorm:"type:text;not null" json:"options"`              // JSON格式存储选项
	Answer       string         `gorm:"type:text;not null" json:"answer"`
	Explanation  string         `gorm:"type:text" json:"explanation,omitempty"`
	Keywords     string         `gorm:"type:VARCHAR(255)" json:"keywords,omitempty"`
	Language     string         `gorm:"type:VARCHAR(50);not null" json:"language"`
	AiModel      string         `gorm:"type:VARCHAR(50);not null" json:"ai_model"`
	UserID       int64          `gorm:"not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"` // 关联用户表
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 显式指定表名
func (Question) TableName() string {
	return "questions"
}
