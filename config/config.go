package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config 结构体映射所有环境变量配置
type Config struct {
	// AI API 配置
	TongyiAPIKey   string // 通义千问 API 密钥
	DeepseekAPIKey string // Deepseek API 密钥

	// 支持的编程语言（解析为切片方便使用）
	SupportedLanguages []string

	// 数据库配置
	DBPath string // 数据库文件路径

	// JWT 配置
	JWTSecret      string // JWT 加密密钥
	JWTExpireHours int    // JWT 过期时间（小时）

	// 服务器配置
	ServerPort int    // 服务器端口
	GinMode    string // Gin 运行模式（debug/release/test）
}

// LoadConfig 加载并合并 .env 和 .env.local 配置
func LoadConfig() (*Config, error) {
	// 先加载基础配置 .env
	if err := godotenv.Load(".env"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("加载 .env 失败: %w", err)
	}

	// 再加载 .env.local 覆盖基础配置
	if err := godotenv.Load(".env.local"); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("加载 .env.local 失败: %w", err)
	}

	// 解析配置项（带默认值处理）
	cfg := &Config{
		// AI API 配置（无默认值，必须在 .env 或 .env.local 中定义）
		TongyiAPIKey:   getEnv("TONGYI_API_KEY", ""),
		DeepseekAPIKey: getEnv("DEEPSEEK_API_KEY", ""),

		// 支持的编程语言（默认空切片，解析为 []string）
		SupportedLanguages: parseLanguages(getEnv("SUPPORTED_LANGUAGES", "")),

		// 数据库配置（默认当前目录下的 exam_system.db）
		DBPath: getEnv("DB_PATH", "./exam_system.db"),

		// JWT 配置（默认密钥为空，过期时间 24 小时）
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpireHours: getEnvAsInt("JWT_EXPIRE_HOURS", 24),

		// 服务器配置（默认端口 8080，Gin 模式为 debug）
		ServerPort: getEnvAsInt("SERVER_PORT", 8080),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}

	// 验证必要的配置项（避免程序启动后因缺失关键配置出错）
	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// 工具函数：获取环境变量，不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// 工具函数：将环境变量解析为 int 类型
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue // 解析失败时返回默认值
	}
	return value
}

// 工具函数：将逗号分隔的字符串解析为编程语言切片
func parseLanguages(langsStr string) []string {
	if langsStr == "" {
		return []string{}
	}
	// 去除空格并分割（处理类似 "Go, Python, Java" 带空格的情况）
	langs := strings.Split(strings.ReplaceAll(langsStr, " ", ""), ",")
	return langs
}

// 验证配置项的合法性
func (c *Config) validate() error {
	// 验证 AI API 密钥（至少需要一个 API 密钥，根据你的业务需求调整）
	if c.TongyiAPIKey == "" && c.DeepseekAPIKey == "" {
		return fmt.Errorf("至少需要配置 TONGYI_API_KEY 或 DEEPSEEK_API_KEY")
	}

	// 验证 JWT 密钥（生产环境必须配置，避免使用默认空值）
	if c.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET 不能为空，请配置密钥")
	}

	// 验证服务器端口（范围 1-65535）
	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("SERVER_PORT 必须在 1-65535 之间，当前值: %d", c.ServerPort)
	}

	// 验证 Gin 模式（只能是 debug/release/test）
	validGinModes := map[string]bool{
		"debug":   true,
		"release": true,
		"test":    true,
	}
	if !validGinModes[c.GinMode] {
		return fmt.Errorf("GIN_MODE 必须是 debug/release/test，当前值: %s", c.GinMode)
	}

	return nil
}
