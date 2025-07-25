# 1、CodeQuizAI 项目介绍

## 功能概述

CodeQuizAI 是一个基于 AI 的编程题库生成与管理系统，主要功能围绕编程相关题目的自动生成、管理和使用展开。系统允许用户通过 AI 模型快速生成特定编程语言的选择题（单选题和多选题），并支持对生成的题目进行编辑、确认和组织成试卷。

核心功能包括：
1. 用户认证（登录/注册）
2. 基于 AI 的编程题目生成（支持多种 AI 模型和编程语言）
3. 临时题目预览与编辑
4. 题目确认与正式入库
5. 试卷管理功能
6. 统计分析功能

## 技术栈

- 后端框架：Golang + Gin
- 数据库：SQLite（通过 gorm 进行 ORM 操作）
- AI 集成：支持通义千问（Tongyi）和 DeepSeek 模型
- 认证机制：JWT（JSON Web Token）
- 代码生成：gorm/gen 用于生成数据访问层代码

## 主要特性

1. **多 AI 模型支持**
    - 集成通义千问和 DeepSeek 两种 AI 模型
    - 可根据配置自动选择可用的 AI 模型生成题目

2. **灵活的题目生成**
    - 支持指定编程语言
    - 可选择题目类型（单选/多选）
    - 支持通过关键词限定题目范围
    - 可指定生成题目数量（1-10道）

3. **完整的题目生命周期管理**
    - 临时存储 AI 生成的题目（temp_questions 表）
    - 支持预览和编辑临时题目
    - 确认后将题目正式入库
    - 记录题目的创建者和相关元数据

4. **安全可靠的系统设计**
    - JWT 身份认证与授权
    - 密码加密存储（SHA-256 哈希）
    - 配置项合法性校验
    - 数据库事务保证数据一致性

5. **可扩展的架构**
    - 分层设计（控制器、服务、数据访问层）
    - 中间件支持（认证、权限控制）
    - 模块化配置管理

6. **便捷的数据库管理**
    - 自动初始化数据库
    - 支持 SQL 脚本执行
    - 软删除功能（通过 gorm.DeletedAt 实现）


# 2、CodeQuizAI 配置说明：环境变量和配置文件

## 配置文件结构

CodeQuizAI 采用 `.env` 文件进行配置管理，支持基础配置与本地配置的分层管理：

1. **基础配置文件**：`.env`  
   存储通用配置，适合纳入版本控制，包含非敏感的基础设置。

2. **本地配置文件**：`.env.local`  
   存储本地环境的个性化配置（如 API 密钥），会覆盖 `.env` 中的同名配置，**不应纳入版本控制**。

加载优先级：`.env.local` > `.env` > 系统环境变量

## 核心环境变量说明

### 1. AI 模型配置（至少配置一个）
```ini
# 通义千问 API 密钥
TONGYI_API_KEY=your_tongyi_api_key_here

# Deepseek API 密钥
DEEPSEEK_API_KEY=your_deepseek_api_key_here
```

### 2. 编程语言支持配置
```ini
# 支持的编程语言（逗号分隔，无空格）
# 示例：SUPPORTED_LANGUAGES=Go,Python,Java,JavaScript,C++,C
SUPPORTED_LANGUAGES=Go,Python,Java,JavaScript,C++,C
```

### 3. 数据库配置
```ini
# SQLite 数据库文件路径（默认：./exam_system.db）
DB_PATH=./CodeQuizAI.db
```

### 4. JWT 认证配置
```ini
# JWT 加密密钥（必填，生产环境需使用强密钥）
JWT_SECRET=your_strong_jwt_secret_key

# JWT 过期时间（小时，默认：24）
JWT_EXPIRE_HOURS=24
```

### 5. 服务器配置
```ini
# 服务器端口（默认：8080）
SERVER_PORT=8080

# Gin 运行模式（debug/release/test，默认：debug）
GIN_MODE=debug
```

## 配置加载与验证流程

1. **加载过程**：
    - 先加载 `.env` 基础配置
    - 再加载 `.env.local` 覆盖基础配置
    - 解析环境变量并转换为对应类型（字符串/整数）

2. **自动验证**：
    - 至少需要一个 AI API 密钥（TONGYI_API_KEY 或 DEEPSEEK_API_KEY）
    - JWT_SECRET 不能为空（确保认证安全）
    - 服务器端口必须在 1-65535 范围内
    - Gin 模式必须是 debug/release/test 中的一种

## 配置示例

**.env 文件示例**：
```ini
# 基础配置
SUPPORTED_LANGUAGES=Go,Python,Java,JavaScript,C++,C
DB_PATH=./CodeQuizAI.db
SERVER_PORT=8080
GIN_MODE=debug
JWT_EXPIRE_HOURS=24
```

**.env.local 文件示例**：
```ini
# 本地敏感配置（覆盖.env）
TONGYI_API_KEY=sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
DEEPSEEK_API_KEY=sk-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy
JWT_SECRET=your_very_strong_secret_key_here
SERVER_PORT=8088
```

## 注意事项
`.env.local` 已包含在 `.gitignore` 中，确保敏感信息不会提交到代码库

配置修改后需重启服务才能生效。


# 3、CodeQuizAI 数据库自动初始化机制说明

## 初始化流程概述

CodeQuizAI 采用自动检测与初始化机制，确保数据库在首次运行时能自动创建并完成基础配置。整个流程由 `main.go` 中的 `initDatabase` 函数主导，核心逻辑是检测数据库文件是否存在，若不存在则创建并执行初始化脚本。

## 关键实现代码

```go
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
```

## 初始化步骤详解

1. **存在性检测**
    - 通过 `os.Stat(dbPath)` 检查指定路径的 SQLite 数据库文件是否存在
    - 若文件不存在（`os.IsNotExist(err)`），则触发初始化流程
    - 若文件已存在，则直接跳过初始化

2. **目录创建**
    - 使用 `os.MkdirAll` 创建数据库文件所在的目录（支持多级目录）
    - 权限设置为 `0755`（所有者可读写执行，组和其他用户可读执行）

3. **数据库文件创建**
    - 通过 `gorm.Open(sqlite.Open(dbPath))` 自动创建空的 SQLite 数据库文件
    - SQLite 特性：连接不存在的文件时会自动创建该文件

4. **初始化脚本执行**
    - 读取 `./migrations/init.sql` 脚本文件内容
    - 使用事务（`tx.Begin()`）确保脚本执行的原子性
    - 按分号分割 SQL 语句并逐条执行
    - 执行成功则提交事务，失败则回滚

## 初始化脚本内容（init.sql）

脚本定义了系统所需的全部数据表结构及默认数据：

1. **核心表结构**
    - `users`：用户表（存储用户账号信息）
    - `questions`：正式题目表（存储确认入库的题目）
    - `papers`：试卷表（试卷基本信息）
    - `paper_questions`：试卷题目关联表（多对多关系）
    - `temp_questions`：临时题目表（存储AI生成的未确认题目）

2. **默认数据**
    - 插入默认管理员账户（`admin`）
    - 明文密码：`123456`
    - 密码哈希：`8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92`

## 配置与路径说明

- 数据库文件路径由配置项 `DB_PATH` 指定（默认：`./CodeQuizAI.db`）
- 初始化脚本固定路径：`./migrations/init.sql`
- 可通过修改环境变量 `DB_PATH` 自定义数据库存储位置

## 注意事项

1. 首次启动时会自动执行初始化，后续启动因文件已存在将跳过
2. 若需重新初始化，需手动删除数据库文件后重启服务
3. 数据库文件不会提交到仓库


# 4、CodeQuizAI 数据库结构
详细的表结构设计文档放置路径:`./migrations/README.md`


# 5、CodeQuizAI 项目 API 文档
正在完善中...