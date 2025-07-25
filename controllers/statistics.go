package controllers

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/services"
	"CodeQuizAI/utils"
	"github.com/gin-gonic/gin"
	"strconv"
)

// GetUserStatistics 获取用户统计信息
func GetUserStatistics(c *gin.Context) {
	// 1. 解析路径参数中的用户ID
	userIDStr := c.Param("id")
	targetUserID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		utils.SendResponse(c, 400, "无效的用户ID", nil)
		return
	}

	// 2. 获取当前登录用户ID（用于权限验证）
	currentUserID, exists := c.Get("user_id")
	if !exists {
		utils.SendResponse(c, 401, "未获取到用户信息", nil)
		return
	}

	// 3. 权限验证：仅允许管理员或用户本人查看统计信息
	currentUserIDInt64 := int64(currentUserID.(uint))
	if currentUserIDInt64 != targetUserID {
		// 检查当前用户是否为管理员
		user, err := dao.Q.User.WithContext(c).Where(dao.User.ID.Eq(currentUserIDInt64)).First()
		if err != nil || user.Role != "admin" {
			utils.SendResponse(c, 403, "无权查看该用户的统计信息", nil)
			return
		}
	}

	// 4. 调用服务层获取统计数据
	statistics, err := services.GetUserStatistics(c.Request.Context(), targetUserID)
	if err != nil {
		utils.SendResponse(c, 500, "获取统计信息失败："+err.Error(), nil)
		return
	}

	// 5. 返回响应
	utils.SendResponse(c, 200, "获取用户统计信息成功", statistics)
}

// GetStatisticsOverview 处理获取系统统计概览的请求
func GetStatisticsOverview(c *gin.Context) {
	// 调用服务层获取统计数据
	overview, err := services.GetStatisticsOverview(c.Request.Context())
	if err != nil {
		utils.SendResponse(c, 500, "获取统计数据失败："+err.Error(), nil)
		return
	}

	utils.SendResponse(c, 200, "获取统计数据成功", overview)
}
