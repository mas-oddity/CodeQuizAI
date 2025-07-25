package models

import (
	"time"
)

// PaperQuestion 对应数据库中的 paper_questions 表
type PaperQuestion struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	PaperID       int64     `gorm:"not null" json:"paper_id"`
	QuestionID    int64     `gorm:"not null" json:"question_id"`
	QuestionOrder int       `gorm:"not null" json:"question_order"` // 题目顺序
	Score         int       `gorm:"default:5" json:"score"`         // 该题分值
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	// 关联模型
	Paper    Paper    `gorm:"foreignKey:PaperID" json:"paper,omitempty"`
	Question Question `gorm:"foreignKey:QuestionID" json:"question,omitempty"`
}

// TableName 显式指定表名
func (PaperQuestion) TableName() string {
	return "paper_questions"
}
