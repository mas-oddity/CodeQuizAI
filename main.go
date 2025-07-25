package main

import (
	"CodeQuizAI/config"
	"CodeQuizAI/dao"
	"CodeQuizAI/middlewares"
	"CodeQuizAI/router"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	// 设置 Gin 运行模式（从配置中读取）
	gin.SetMode(cfg.GinMode)

	// 配置参数
	dbPath := cfg.DBPath
	initScriptPath := "./migrations/init.sql"

	// 初始化数据库
	err = initDatabase(dbPath, initScriptPath)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %w", err)
	}

	// 设置DAO默认数据库连接
	dao.SetDefault(db)

	// 初始化JWT配置
	middlewares.InitJWT(cfg.JWTSecret, time.Duration(cfg.JWTExpireHours)*time.Hour) // 密钥和过期时间

	// 启动Web服务
	log.Println("启动Web服务...")
	r := router.InitRouter()
	if err := r.Run(fmt.Sprintf(":%d", cfg.ServerPort)); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

// 初始化数据库，如果不存在则创建并执行初始化脚本
func initDatabase(dbPath, initScriptPath string) error {
	// 检查数据库文件是否存在
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Println("数据库文件不存在，开始初始化...")

		// 创建数据库文件所在目录
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return fmt.Errorf("创建数据库目录失败: %w", err)
		}

		// 连接数据库（自动创建空文件）
		db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
		if err != nil {
			return fmt.Errorf("连接数据库失败: %w", err)
		}

		// 获取底层 sql.DB 实例
		sqlDB, err := db.DB()
		if err != nil {
			return fmt.Errorf("获取SQL数据库实例失败: %w", err)
		}
		defer sqlDB.Close()

		// 执行初始化脚本
		if err := executeSQLScript(sqlDB, initScriptPath); err != nil {
			return fmt.Errorf("执行初始化脚本失败: %w", err)
		}

		log.Println("数据库初始化完成")

		return nil
	}

	log.Println("数据库文件已存在，跳过初始化")

	return nil
}

// 执行SQL脚本文件
func executeSQLScript(db *sql.DB, scriptPath string) error {
	// 读取SQL脚本
	script, err := os.ReadFile(scriptPath)
	if err != nil {
		return fmt.Errorf("读取SQL脚本失败: %w", err)
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// 按分号分割SQL语句并执行
	for _, stmt := range strings.Split(string(script), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		// 执行SQL语句
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("执行SQL语句失败: %w\n语句: %s", err, stmt)
		}
	}

	return nil
}
