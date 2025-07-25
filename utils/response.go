package utils

import (
	"github.com/gin-gonic/gin"
)

// Response 统一响应格式结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SendResponse 发送统一格式的响应
func SendResponse(c *gin.Context, code int, message string, data interface{}) {
	resp := Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
	c.JSON(code, resp)
}
