package middlewares

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTClaims 自定义 JWT 声明
type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.StandardClaims
}

// JWTConfig 配置信息
type JWTConfig struct {
	SecretKey       string        // 密钥
	TokenExpireTime time.Duration // 令牌过期时间（小时）
}

var jwtConfig JWTConfig

// InitJWT 初始化 JWT 配置
func InitJWT(secretKey string, expireTime time.Duration) {
	jwtConfig = JWTConfig{
		SecretKey:       secretKey,
		TokenExpireTime: expireTime,
	}
}

// AuthMiddleware JWT 鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少授权头"})
			c.Abort()
			return
		}

		// 验证授权头格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "授权头格式错误"})
			c.Abort()
			return
		}

		// 解析 token
		tokenStr := parts[1]
		claims := &JWTClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
			}
			return []byte(jwtConfig.SecretKey), nil
		})

		if err != nil {
			if ve, ok := err.(*jwt.ValidationError); ok {
				if ve.Errors&jwt.ValidationErrorExpired != 0 {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "令牌已过期"})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
				}
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "解析令牌失败"})
			}
			c.Abort()
			return
		}

		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的令牌"})
			c.Abort()
			return
		}

		// 将用户ID存储到上下文，供后续处理使用
		c.Set("user_id", claims.UserID)
		c.Next()
	}
}

// GenerateToken 生成 JWT 令牌
func GenerateToken(userID int64) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtConfig.TokenExpireTime * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "codequizai",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtConfig.SecretKey))
}

// ParseToken 解析 JWT 令牌
func ParseToken(tokenStr string) (*JWTClaims, error) {
	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		}
		return []byte(jwtConfig.SecretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("无效的令牌")
	}

	return claims, nil
}
