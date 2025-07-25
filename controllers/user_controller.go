package controllers

import (
	"CodeQuizAI/services"
	"CodeQuizAI/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListUsers 管理员获取用户列表接口
func ListUsers(c *gin.Context) {
	// 1. 解析分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")

	// 2. 调用Service层处理业务
	userList, pagination, err := services.ListUsers(
		c,
		page,
		pageSize,
		username,
	)

	// 3. 处理响应
	if err != nil {
		utils.SendResponse(c, 500, err.Error(), nil)
		return
	}

	utils.SendResponse(c, 200, "获取用户列表成功", gin.H{
		"list":       userList,
		"pagination": pagination,
	})
}

// UpdateUser 允许登录用户更新自己的信息（用户名/密码）
func UpdateUser(c *gin.Context) {
	// 1. 解析路径参数中的目标用户ID
	targetUserID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "无效的用户ID", nil)
		return
	}

	// 2. 获取当前登录用户ID（从JWT中间件存入的上下文）
	currentUserID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, http.StatusUnauthorized, "未获取到用户信息", nil)
		return
	}

	// 3. 校验：仅允许修改自己的信息（目标ID必须等于当前登录用户ID）
	if int64(currentUserID.(uint)) != targetUserID {
		utils.SendResponse(c, http.StatusForbidden, "无权修改他人信息", nil)
		return
	}

	// 4. 解析请求体参数（用户名和密码为可选，但至少传一个）
	var req struct {
		Username string `json:"username" binding:"omitempty,min=3,max=50"` // 可选，3-50字符
		Password string `json:"password" binding:"omitempty,min=6"`        // 可选，至少6位
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendResponse(c, http.StatusBadRequest, "参数校验失败: "+err.Error(), nil)
		return
	}
	// 校验至少传一个字段
	if req.Username == "" && req.Password == "" {
		utils.SendResponse(c, http.StatusBadRequest, "至少需要修改用户名或密码", nil)
		return
	}

	// 5. 调用Service层处理业务
	err = services.UpdateUser(c, targetUserID, req.Username, req.Password)
	if err != nil {
		utils.SendResponse(c, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	// 6. 返回成功响应
	utils.SendResponse(c, http.StatusOK, "用户信息更新成功", nil)
}
