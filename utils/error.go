package utils

import "errors"

// 自定义错误变量
var (
	ErrQuestionNotFound  = errors.New("题目不存在")
	ErrNoPermission      = errors.New("没有权限")
	ErrPaperNotFound     = errors.New("试卷不存在或已被删除")
	ErrDuplicateQuestion = errors.New("题目已存在于试卷中")
	ErrInvalidOrder      = errors.New("题目顺序重复或无效")
	ErrScoreExceedTotal  = errors.New("题目总分超过试卷上限")
)
