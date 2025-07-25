package controllers

import (
	"CodeQuizAI/config"
	"CodeQuizAI/models"
	"CodeQuizAI/services"
	"CodeQuizAI/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"strconv"
)

// GenerateQuestionResponse 生成题目的响应数据
type GenerateQuestionResponse struct {
	PreviewID string                `json:"preview_id"` // 预览批次ID
	Questions []models.TempQuestion `json:"questions"`  // 生成的临时题目
}

// GenerateQuestions 处理题目生成请求
func GenerateQuestions(c *gin.Context) {
	// 1. 解析并验证请求参数
	var req services.GenerateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 2. 获取当前登录用户ID（从JWT中间件上下文）
	userID, _ := c.Get("user_id")
	userIDInt64, _ := userID.(int64)

	// 3. 验证编程语言是否在配置的支持列表中
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}
	if !services.IsLanguageSupported(req.Language, cfg.SupportedLanguages) {
		utils.SendResponse(c, 400, "不支持的编程语言："+req.Language, nil)
		return
	}

	// 4. 调用服务层生成题目
	previewID := uuid.New().String() // 生成预览批次ID
	tempQuestions, err := services.GenerateQuestions(
		c.Request.Context(),
		previewID,
		userIDInt64,
		req,
		cfg,
	)
	if err != nil {
		utils.SendResponse(c, 500, "生成题目失败："+err.Error(), nil)
		return
	}

	// 5. 返回成功响应
	data := GenerateQuestionResponse{
		PreviewID: previewID,
		Questions: tempQuestions,
	}
	utils.SendResponse(c, 200, "题目生成成功", data)
}

// ConfirmQuestions 处理题目确认入库请求
func ConfirmQuestions(c *gin.Context) {
	// 1. 解析请求参数
	var req services.ConfirmQuestionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 2. 获取当前用户ID
	userID, _ := c.Get("user_id")
	userIDInt64, _ := userID.(int64)

	// 3. 调用服务层执行确认逻辑
	result, err := services.ConfirmQuestions(
		c.Request.Context(),
		req.PreviewID,
		req.Selected,
		userIDInt64,
	)
	if err != nil {
		utils.SendResponse(c, 500, "确认题目失败："+err.Error(), nil)
		return
	}

	// 4. 返回成功响应
	utils.SendResponse(c, 200, "题目已成功入库", result)
}

// GetQuestions 查询当前用户的题目列表
func GetQuestions(c *gin.Context) {
	// 1. 解析查询参数
	var req services.GetQuestionsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 2. 设置默认参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 50 {
		req.PageSize = 10
	}
	if req.Sort == "" {
		req.Sort = "created_at_desc" // 默认按创建时间降序
	}

	// 3. 获取当前用户ID
	userID, _ := c.Get("user_id")
	userIDInt64, _ := userID.(int64)

	// 4. 调用服务层查询
	result, err := services.GetUserQuestions(
		c.Request.Context(),
		userIDInt64,
		req,
	)
	if err != nil {
		utils.SendResponse(c, 500, "查询题目失败："+err.Error(), nil)
		return
	}

	// 5. 返回响应
	utils.SendResponse(c, 200, "查询成功", result)
}

// UpdateQuestion 更新指定题目
func UpdateQuestion(c *gin.Context) {
	// 1. 解析路径参数（题目ID）
	idStr := c.Param("id")
	questionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的题目ID", nil)
		return
	}

	// 2. 解析请求体
	var req services.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 3. 验证至少传递一个更新字段
	if !HasAnyField(req) {
		utils.SendResponse(c, 400, "至少需要传递一个更新字段", nil)
		return
	}

	// 4. 获取当前用户ID
	userID, _ := c.Get("user_id")
	userIDInt64, _ := userID.(int64)

	// 5. 调用服务层执行更新
	result, err := services.UpdateQuestion(
		c.Request.Context(),
		questionID,
		userIDInt64,
		req,
	)
	if err != nil {
		// 区分不同错误类型
		if errors.Is(err, utils.ErrQuestionNotFound) {
			utils.SendResponse(c, 404, "题目不存在或已被删除", nil)
		} else if errors.Is(err, utils.ErrNoPermission) {
			utils.SendResponse(c, 403, "没有权限更新该题目", nil)
		} else {
			utils.SendResponse(c, 500, "更新题目失败："+err.Error(), nil)
		}
		return
	}

	// 6. 返回成功响应
	utils.SendResponse(c, 200, "题目更新成功", result)
}

// HasAnyField 检查是否至少传递了一个字段
func HasAnyField(u services.UpdateQuestionRequest) bool {
	return u.Title != "" ||
		u.QuestionType != "" ||
		u.Options != "" ||
		u.Answer != "" ||
		u.Explanation != "" ||
		u.Keywords != ""
}

// DeleteQuestion 软删除指定题目
func DeleteQuestion(c *gin.Context) {
	// 1. 解析路径参数（题目ID）
	idStr := c.Param("id")
	questionID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的题目ID", nil)
		return
	}

	// 2. 获取当前用户ID
	userID, _ := c.Get("user_id")
	userIDInt64, _ := userID.(int64)

	// 3. 调用服务层执行删除
	result, err := services.DeleteQuestion(
		c.Request.Context(),
		questionID,
		userIDInt64,
	)
	if err != nil {
		// 区分错误类型
		if errors.Is(err, utils.ErrQuestionNotFound) {
			utils.SendResponse(c, 404, "题目不存在或已被删除", nil)
		} else if errors.Is(err, utils.ErrNoPermission) {
			utils.SendResponse(c, 403, "没有权限删除该题目", nil)
		} else {
			utils.SendResponse(c, 500, "删除题目失败："+err.Error(), nil)
		}
		return
	}

	// 4. 返回成功响应
	utils.SendResponse(c, 200, "题目已成功删除", result)
}
