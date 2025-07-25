-- 创建用户表
CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user',
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL
    );

-- 创建题目表
CREATE TABLE IF NOT EXISTS questions (
                                         id INTEGER PRIMARY KEY AUTOINCREMENT,
                                         title TEXT NOT NULL,
                                         question_type VARCHAR(20) NOT NULL,  -- 'single' 或 'multiple'
    options TEXT NOT NULL,               -- JSON格式存储选项
    answer TEXT NOT NULL,
    explanation TEXT,
    keywords VARCHAR(255),
    language VARCHAR(50) NOT NULL,       -- 编程语言
    ai_model VARCHAR(50) NOT NULL,       -- 使用的AI模型
    user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (user_id) REFERENCES users(id)
    );

-- 创建试卷表
CREATE TABLE IF NOT EXISTS papers (
                                      id INTEGER PRIMARY KEY AUTOINCREMENT,
                                      title VARCHAR(255) NOT NULL,
    description TEXT,
    total_score INTEGER DEFAULT 100,
    creator_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    FOREIGN KEY (creator_id) REFERENCES users(id)
    );

-- 创建试卷题目关联表
CREATE TABLE IF NOT EXISTS paper_questions (
                                               id INTEGER PRIMARY KEY AUTOINCREMENT,
                                               paper_id INTEGER NOT NULL,
                                               question_id INTEGER NOT NULL,
                                               question_order INTEGER NOT NULL,     -- 题目顺序
                                               score INTEGER DEFAULT 5,             -- 该题分值
                                               created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
                                               FOREIGN KEY (paper_id) REFERENCES papers(id),
    FOREIGN KEY (question_id) REFERENCES questions(id)
    );

-- 创建临时临时题目表（用于存储未确认入库的AI生成题目）
CREATE TABLE IF NOT EXISTS temp_questions (
                                              id INTEGER PRIMARY KEY AUTOINCREMENT,
                                              preview_id VARCHAR(64) NOT NULL, -- 预览批次唯一标识（UUID）
    temp_id VARCHAR(64) NOT NULL,    -- 单题临时ID（preview_id + 序号）
    title TEXT NOT NULL,             -- 题目标题
    question_type VARCHAR(20) NOT NULL, -- 题目类型（single/multiple）
    options TEXT NOT NULL,           -- 选项（JSON格式）
    answer TEXT NOT NULL,            -- 答案
    explanation TEXT,                -- 解析（可选）
    keywords VARCHAR(255),           -- 关键词（可选）
    language VARCHAR(50) NOT NULL,   -- 编程语言
    ai_model VARCHAR(50) NOT NULL,   -- 使用的AI模型
    user_id INTEGER NOT NULL,        -- 关联用户ID
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME NULL,
    UNIQUE(preview_id, temp_id, user_id),
    FOREIGN KEY (user_id) REFERENCES users(id)
    );

-- 插入默认管理员账户
INSERT OR IGNORE INTO users (username, password_hash, role)
VALUES ('admin', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92', 'admin');