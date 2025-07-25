package middlewares

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

// AdminMiddleware 管理员权限校验中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取 user_id（由 AuthMiddleware 设置）
		userID, exists := c.Get("user_id")
		if !exists {
			utils.SendResponse(c, http.StatusUnauthorized, "未获取到用户信息", nil)
			c.Abort()
			return
		}

		// 转换 user_id 类型（注意：models.User.ID 是 int64 类型，需与 JWT 中存储的类型一致）
		// 从 models/user.go 可知 ID 为 int64，而 JWTClaims 中 UserID 是 uint，此处需统一类型
		uid, ok := userID.(int64)
		if !ok {
			utils.SendResponse(c, http.StatusInternalServerError, "用户ID类型错误", nil)
			c.Abort()
			return
		}

		// 查询用户角色（通过 DAO 层获取用户信息）
		user, err := dao.Q.User.WithContext(c).Where(dao.User.ID.Eq(uid)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				utils.SendResponse(c, http.StatusNotFound, "用户不存在", nil)
			} else {
				utils.SendResponse(c, http.StatusInternalServerError, "查询用户信息失败", nil)
			}
			c.Abort()
			return
		}

		// 校验角色是否为管理员
		if user.Role != "admin" {
			utils.SendResponse(c, http.StatusForbidden, "权限不足，仅管理员可访问", nil)
			c.Abort()
			return
		}

		// 校验通过，继续处理请求
		c.Next()
	}
}
