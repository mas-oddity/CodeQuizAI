package controllers

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/middlewares"
	"CodeQuizAI/models"
	"CodeQuizAI/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

// Register 处理用户注册请求（用户名和密码）
func Register(c *gin.Context) {
	// 定义注册请求结构体，包含用户名和密码
	var registerRequest struct {
		Username string `json:"username" binding:"required,min=3,max=20"` // 用户名3-20个字符
		Password string `json:"password" binding:"required,min=6,max=32"` // 密码6-32个字符
	}

	// 解析并验证请求体
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		utils.SendResponse(c, 400, "无效的请求参数: "+err.Error(), nil)
		return
	}

	// 验证用户名格式（仅允许字母、数字和下划线）
	if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(registerRequest.Username) {
		utils.SendResponse(c, 400, "用户名只能包含字母、数字和下划线", nil)
		return
	}

	// 检查用户名是否已存在
	userDAO := dao.Q.User
	existingUser, err := userDAO.WithContext(c).
		Where(dao.User.Username.Eq(registerRequest.Username)).
		First()

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.SendResponse(c, 500, "数据库查询出错", nil)
		return
	}

	// 如果用户已存在
	if existingUser != nil {
		utils.SendResponse(c, 400, "用户名已被注册", nil)
		return
	}

	// 计算 SHA-256 哈希
	hash := sha256.Sum256([]byte(registerRequest.Password))
	// 转为十六进制字符串
	hashStr := hex.EncodeToString(hash[:])

	// 创建新用户
	newUser := &models.User{
		Username:     registerRequest.Username,
		PasswordHash: hashStr,
	}

	dao.User.WithContext(c).Create(newUser)

	// 为新用户生成令牌
	token, err := middlewares.GenerateToken(newUser.ID)
	if err != nil {
		utils.SendResponse(c, 500, "生成令牌出错", nil)
		return
	}

	// 返回注册成功响应
	data := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"username": registerRequest.Username,
		},
	}
	utils.SendResponse(c, 201, "注册成功", data)
}
