package controllers

import (
	"CodeQuizAI/services"
	"CodeQuizAI/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// GetPapers 查询当前用户的试卷列表
func GetPapers(c *gin.Context) {
	// 1. 解析查询参数
	var req services.GetPapersRequest
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

	// 3. 获取当前用户ID（作为试卷创建者ID）
	userID, _ := c.Get("user_id")
	creatorID, _ := userID.(int64)

	// 4. 调用服务层查询
	result, err := services.GetUserPapers(
		c.Request.Context(),
		creatorID,
		req,
	)
	if err != nil {
		utils.SendResponse(c, 500, "查询试卷失败："+err.Error(), nil)
	}

	// 5. 返回响应
	utils.SendResponse(c, 200, "查询试卷列表成功", result)
}

// CreatePaper 创建新试卷
func CreatePaper(c *gin.Context) {
	// 1. 解析请求体参数
	var req services.CreatePaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 2. 获取当前用户ID（作为创建者ID）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, 401, "未获取到用户信息", nil)
		return
	}
	creatorID, _ := userID.(int64)

	// 3. 调用服务层创建试卷
	paper, err := services.CreatePaper(
		c.Request.Context(),
		creatorID,
		req,
	)
	if err != nil {
		utils.SendResponse(c, 500, "创建试卷失败："+err.Error(), nil)
		return
	}

	// 4. 返回成功响应
	utils.SendResponse(c, 201, "试卷创建成功", paper)
}

// GetPaperDetail 获取试卷详情（包含关联题目）
func GetPaperDetail(c *gin.Context) {
	// 1. 解析路径参数（试卷ID）
	idStr := c.Param("id")
	paperID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	// 2. 获取当前用户ID（验证权限）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, 401, "未获取到用户信息", nil)
		return
	}
	creatorID, _ := userID.(int64)

	// 3. 调用服务层查询详情
	paperDetail, err := services.GetPaperDetail(
		c.Request.Context(),
		paperID,
		creatorID,
	)
	if err != nil {
		if errors.Is(err, utils.ErrPaperNotFound) {
			utils.SendResponse(c, 404, "试卷不存在或已被删除", nil)
		} else if errors.Is(err, utils.ErrNoPermission) {
			utils.SendResponse(c, 403, "没有权限查看该试卷", nil)
		} else {
			utils.SendResponse(c, 500, "查询试卷详情失败："+err.Error(), nil)
		}
		return
	}

	// 4. 返回响应
	utils.SendResponse(c, 200, "查询试卷详情成功", paperDetail)
}

// DeletePaper 软删除试卷（同时删除试卷与题目的关联关系）
func DeletePaper(c *gin.Context) {
	// 1. 解析路径参数（试卷ID）
	idStr := c.Param("id")
	paperID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	// 2. 获取当前用户ID（从JWT中间件）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, 401, "未获取到用户信息", nil)
		return
	}
	creatorID, _ := userID.(int64)

	// 3. 调用服务层执行删除
	result, err := services.DeletePaper(
		c.Request.Context(),
		paperID,
		creatorID,
	)
	if err != nil {
		if errors.Is(err, utils.ErrPaperNotFound) {
			utils.SendResponse(c, 404, "试卷不存在或已被删除", nil)
		} else if errors.Is(err, utils.ErrNoPermission) {
			utils.SendResponse(c, 403, "没有权限删除该试卷", nil)
		} else {
			utils.SendResponse(c, 500, "删除试卷失败："+err.Error(), nil)
		}
		return
	}

	// 4. 返回成功响应
	utils.SendResponse(c, 200, "试卷已成功删除", result)
}

// AddQuestionsToPaper 处理向试卷添加题目的请求
func AddQuestionsToPaper(c *gin.Context) {
	// 1. 解析路径参数（试卷ID）
	paperIDStr := c.Param("id")
	paperID, err := strconv.ParseInt(paperIDStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	// 2. 解析请求体
	var req services.AddQuestionsToPaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误: "+err.Error(), nil)
		return
	}

	// 3. 验证至少添加一个题目
	if len(req.Items) == 0 {
		utils.SendResponse(c, 400, "至少需要添加一个题目", nil)
		return
	}

	// 4. 获取当前用户ID（试卷创建者）
	userID, _ := c.Get("user_id")
	creatorID, _ := userID.(int64)

	// 5. 调用服务层执行添加
	result, err := services.AddQuestionsToPaper(
		c.Request.Context(),
		paperID,
		creatorID,
		req,
	)
	if err != nil {
		// 区分错误类型
		switch {
		case errors.Is(err, utils.ErrPaperNotFound):
			utils.SendResponse(c, 404, "试卷不存在或无权限", nil)
		case errors.Is(err, utils.ErrInvalidOrder):
			utils.SendResponse(c, 400, "题目顺序重复或无效", nil)
		case strings.Contains(err.Error(), "题目总分超过试卷上限"):
			utils.SendResponse(c, 400, err.Error(), nil)
		default:
			utils.SendResponse(c, 500, "添加题目失败: "+err.Error(), nil)
		}
		return
	}

	// 6. 返回成功响应
	utils.SendResponse(c, 200, "题目添加成功", result)
}

// RemoveQuestionFromPaper 从试卷中移除题目
func RemoveQuestionFromPaper(c *gin.Context) {
	// 解析路径参数
	paperID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	questionID, err := strconv.ParseInt(c.Param("questionID"), 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的题目ID", nil)
		return
	}

	// 获取当前用户ID
	userID, _ := c.Get("user_id")
	creatorID, _ := userID.(int64)

	// 调用服务层
	result, err := services.RemoveQuestionFromPaper(
		c.Request.Context(),
		paperID,
		questionID,
		creatorID,
	)
	if err != nil {
		if errors.Is(err, utils.ErrPaperNotFound) {
			utils.SendResponse(c, 404, "试卷不存在", nil)
		} else {
			utils.SendResponse(c, 500, "移除题目失败："+err.Error(), nil)
		}
		return
	}

	utils.SendResponse(c, 200, "题目移除成功", result)
}

// UpdateQuestionOrder 调整试卷中题目的顺序
func UpdateQuestionOrder(c *gin.Context) {
	// 1. 解析路径参数（试卷ID）
	paperID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	// 2. 解析请求体
	var req services.UpdateQuestionOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 3. 验证请求参数不为空
	if len(req.QuestionOrders) == 0 {
		utils.SendResponse(c, 400, "至少需要传递一个题目顺序信息", nil)
		return
	}

	// 4. 获取当前用户ID
	userID, _ := c.Get("user_id")
	creatorID, _ := userID.(int64)

	// 5. 调用服务层执行更新
	result, err := services.UpdateQuestionOrder(
		c.Request.Context(),
		paperID,
		creatorID,
		req,
	)
	if err != nil {
		if errors.Is(err, utils.ErrPaperNotFound) {
			utils.SendResponse(c, 404, "试卷不存在", nil)
		} else {
			utils.SendResponse(c, 500, "调整题目顺序失败："+err.Error(), nil)
		}
		return
	}

	// 6. 返回响应
	utils.SendResponse(c, 200, "题目顺序调整成功", result)
}

// UpdatePaper 更新试卷信息
func UpdatePaper(c *gin.Context) {
	// 1. 解析路径参数（试卷ID）
	paperID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的试卷ID", nil)
		return
	}

	// 2. 解析请求体参数
	var req services.UpdatePaperRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, 400, "参数错误："+err.Error(), nil)
		return
	}

	// 3. 验证至少传递一个更新字段
	if req.Title == "" && req.Description == "" && req.TotalScore <= 0 {
		utils.SendResponse(c, 400, "至少需要传递一个更新字段（标题/描述/总分）", nil)
		return
	}

	// 4. 获取当前用户ID（作为创建者ID验证权限）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, 401, "未获取到用户信息", nil)
		return
	}
	creatorID, _ := userID.(int64)

	// 5. 调用服务层执行更新
	result, err := services.UpdatePaper(
		c.Request.Context(),
		paperID,
		creatorID,
		req,
	)
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrPaperNotFound):
			utils.SendResponse(c, 404, "试卷不存在", nil)
		case errors.Is(err, utils.ErrNoPermission):
			utils.SendResponse(c, 403, "没有权限更新该试卷", nil)
		default:
			utils.SendResponse(c, 500, "更新试卷失败："+err.Error(), nil)
		}
		return
	}

	// 6. 返回成功响应
	utils.SendResponse(c, 200, "试卷更新成功", result)
}
