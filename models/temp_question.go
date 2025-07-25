package models

import (
	"time"

	"gorm.io/gorm"
)

// TempQuestion 对应数据库中的 temp_questions 表（临时存储AI生成的未确认题目）
type TempQuestion struct {
	ID           int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	PreviewID    string         `gorm:"type:VARCHAR(64);not null" json:"preview_id"`    // 预览批次ID（UUID）
	TempID       string         `gorm:"type:VARCHAR(64);not null" json:"temp_id"`       // 单题临时ID
	Title        string         `gorm:"type:text;not null" json:"title"`                // 题目标题
	QuestionType string         `gorm:"type:VARCHAR(20);not null" json:"question_type"` // 题目类型（single/multiple）
	Options      string         `gorm:"type:text;not null" json:"options"`              // 选项（JSON格式字符串）
	Answer       string         `gorm:"type:text;not null" json:"answer"`               // 答案
	Explanation  string         `gorm:"type:text" json:"explanation,omitempty"`         // 解析（可选）
	Keywords     string         `gorm:"type:VARCHAR(255)" json:"keywords,omitempty"`    // 关键词（可选）
	Language     string         `gorm:"type:VARCHAR(50);not null" json:"language"`      // 编程语言
	AiModel      string         `gorm:"type:VARCHAR(50);not null" json:"ai_model"`      // 使用的AI模型
	UserID       int64          `gorm:"not null" json:"user_id"`                        // 关联用户ID
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`               // 创建时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`              // 软删除字段
}

// TableName 显式指定表名
func (TempQuestion) TableName() string {
	return "temp_questions"
}
