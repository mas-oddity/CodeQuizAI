package controllers

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/middlewares"
	"CodeQuizAI/utils"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Login 处理用户登录请求
func Login(c *gin.Context) {
	// 定义一个结构体来接收请求体中的用户名和密码
	var loginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// 从请求体中解析用户名和密码
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		utils.SendResponse(c, 400, "无效的请求体", nil)
		return
	}

	// 从数据库中查找用户
	userDAO := dao.Q.User
	user, err := userDAO.WithContext(c).Where(dao.User.Username.Eq(loginRequest.Username)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SendResponse(c, 400, "用户名或密码错误", nil)
		} else {
			utils.SendResponse(c, 500, "数据库查询出错", nil)
		}
		return
	}

	// 计算 SHA-256 哈希
	hash := sha256.Sum256([]byte(loginRequest.Password))
	// 转为十六进制字符串
	hashStr := hex.EncodeToString(hash[:])
	if hashStr != user.PasswordHash {
		utils.SendResponse(c, 400, "用户名或密码错误", nil)
		return
	}

	// 生成 JWT 令牌
	token, err := middlewares.GenerateToken((user.ID))
	if err != nil {
		utils.SendResponse(c, 500, "生成令牌出错", nil)
		return
	}

	// 返回成功响应
	data := map[string]string{
		"token": token,
	}
	utils.SendResponse(c, 200, "登录成功", data)
}
