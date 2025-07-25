# 数据库表结构文档

## 1. users 表
### 用途说明
存储系统用户信息，包括管理员和普通用户，用于身份认证和权限管理。

### 字段列表
| 字段名           | 类型         | 说明                          |
|------------------|--------------|-------------------------------|
| id               | INTEGER      | 主键，自增                     |
| username         | VARCHAR(50)  | 用户名，唯一且非空             |
| password_hash    | VARCHAR(255) | 密码哈希值，非空               |
| role             | VARCHAR(20)  | 角色，默认值为'user'           |
| created_at       | DATETIME     | 创建时间，默认当前时间戳       |
| updated_at       | DATETIME     | 更新时间，默认当前时间戳       |
| deleted_at       | DATETIME     | 软删除标记，为空表示未删除     |

### 索引和约束
- 主键约束：`id` 为主键
- 唯一约束：`username` 字段唯一
- 默认值约束：`role` 默认为 'user'
- 非空约束：`username`、`password_hash` 为非空字段

### 关联关系
- 被 `questions` 表关联（一对多）
- 被 `temp_questions` 表关联（一对多）
- 被 `papers` 表关联（一对多）


## 2. questions 表
### 用途说明
存储系统中的正式题目信息，包括选择题的题干、选项、答案等内容。

### 字段列表
| 字段名           | 类型         | 说明                          |
|------------------|--------------|-------------------------------|
| id               | INTEGER      | 主键，自增                     |
| title            | TEXT         | 题目标题，非空                 |
| question_type    | VARCHAR(20)  | 题目类型（'single'或'multiple'），非空 |
| options          | TEXT         | 选项，JSON格式存储，非空       |
| answer           | TEXT         | 答案，非空                     |
| explanation      | TEXT         | 解析，可选                     |
| keywords         | VARCHAR(255) | 关键词，可选                   |
| language         | VARCHAR(50)  | 编程语言，非空                 |
| ai_model         | VARCHAR(50)  | 使用的AI模型，非空             |
| user_id          | INTEGER      | 创建者ID，非空                 |
| created_at       | DATETIME     | 创建时间，默认当前时间戳       |
| updated_at       | DATETIME     | 更新时间，默认当前时间戳       |
| deleted_at       | DATETIME     | 软删除标记，为空表示未删除     |

### 索引和约束
- 主键约束：`id` 为主键
- 非空约束：`title`、`question_type`、`options`、`answer`、`language`、`ai_model`、`user_id` 为非空字段
- 外键约束：`user_id` 关联 `users.id`

### 关联关系
- 关联 `users` 表（多对一）：`user_id` → `users.id`
- 被 `paper_questions` 表关联（一对多）


## 3. papers 表
### 用途说明
存储试卷信息，包括试卷标题、描述、总分等。

### 字段列表
| 字段名           | 类型         | 说明                          |
|------------------|--------------|-------------------------------|
| id               | INTEGER      | 主键，自增                     |
| title            | VARCHAR(255) | 试卷标题，非空                 |
| description      | TEXT         | 试卷描述，可选                 |
| total_score      | INTEGER      | 试卷总分，默认值为100          |
| creator_id       | INTEGER      | 创建者ID，非空                 |
| created_at       | DATETIME     | 创建时间，默认当前时间戳       |
| updated_at       | DATETIME     | 更新时间，默认当前时间戳       |
| deleted_at       | DATETIME     | 软删除标记，为空表示未删除     |

### 索引和约束
- 主键约束：`id` 为主键
- 非空约束：`title`、`creator_id` 为非空字段
- 默认值约束：`total_score` 默认为100
- 外键约束：`creator_id` 关联 `users.id`

### 关联关系
- 关联 `users` 表（多对一）：`creator_id` → `users.id`
- 被 `paper_questions` 表关联（一对多）


## 4. paper_questions 表
### 用途说明
试卷与题目的关联表，记录试卷包含的题目、题目顺序和分值。

### 字段列表
| 字段名           | 类型         | 说明                          |
|------------------|--------------|-------------------------------|
| id               | INTEGER      | 主键，自增                     |
| paper_id         | INTEGER      | 试卷ID，非空                   |
| question_id      | INTEGER      | 题目ID，非空                   |
| question_order   | INTEGER      | 题目顺序，非空                 |
| score            | INTEGER      | 题目分值，默认值为5            |
| created_at       | DATETIME     | 创建时间，默认当前时间戳       |

### 索引和约束
- 主键约束：`id` 为主键
- 非空约束：`paper_id`、`question_id`、`question_order` 为非空字段
- 默认值约束：`score` 默认为5
- 外键约束：`paper_id` 关联 `papers.id`，`question_id` 关联 `questions.id`

### 关联关系
- 关联 `papers` 表（多对一）：`paper_id` → `papers.id`
- 关联 `questions` 表（多对一）：`question_id` → `questions.id`


## 5. temp_questions 表
### 用途说明
存储未确认入库的AI生成临时题目，用于预览和确认。

### 字段列表
| 字段名           | 类型         | 说明                          |
|------------------|--------------|-------------------------------|
| id               | INTEGER      | 主键，自增                     |
| preview_id       | VARCHAR(64)  | 预览批次唯一标识（UUID），非空 |
| temp_id          | VARCHAR(64)  | 单题临时ID，非空               |
| title            | TEXT         | 题目标题，非空                 |
| question_type    | VARCHAR(20)  | 题目类型（'single'或'multiple'），非空 |
| options          | TEXT         | 选项，JSON格式存储，非空       |
| answer           | TEXT         | 答案，非空                     |
| explanation      | TEXT         | 解析，可选                     |
| keywords         | VARCHAR(255) | 关键词，可选                   |
| language         | VARCHAR(50)  | 编程语言，非空                 |
| ai_model         | VARCHAR(50)  | 使用的AI模型，非空             |
| user_id          | INTEGER      | 关联用户ID，非空               |
| created_at       | DATETIME     | 创建时间，默认当前时间戳       |
| deleted_at       | DATETIME     | 软删除标记，为空表示未删除     |

### 索引和约束
- 主键约束：`id` 为主键
- 唯一约束：`(preview_id, temp_id, user_id)` 组合唯一
- 非空约束：`preview_id`、`temp_id`、`title`、`question_type`、`options`、`answer`、`language`、`ai_model`、`user_id` 为非空字段
- 外键约束：`user_id` 关联 `users.id`

### 关联关系
- 关联 `users` 表（多对一）：`user_id` → `users.id`


## 表关联关系图
```
+-------------+       +---------------+       +------------------+
|   users     |       |   questions   |       |  paper_questions |
+-------------+       +---------------+       +------------------+
| id          |<----->| id            |<----->| id               |
| username    |       | title         |       | paper_id         |
| password_hash|      | question_type |       | question_id      |
| role        |       | options       |       | question_order   |
| created_at  |       | answer        |       | score            |
| updated_at  |       | explanation   |       | created_at       |
| deleted_at  |       | keywords      |       +------------------+
+-------------+       | language      |              ^
       ^              | ai_model      |              |
       |              | user_id       |              |
       |              | created_at    |              |
       |              | updated_at    |              |
       |              | deleted_at    |              |
       |              +---------------+              |
       |                                             |
       |              +---------------+              |
       |              |    papers     |              |
       |              +---------------+              |
       +------------->| id            |<------------>|
                      | title         |              |
                      | description   |              |
                      | total_score   |              |
                      | creator_id    |              |
                      | created_at    |              |
                      | updated_at    |              |
                      | deleted_at    |              |
                      +---------------+              |
                                                     |
+-------------+                                      |
| temp_questions|                                     |
+-------------+                                      |
| id          |                                      |
| preview_id  |                                      |
| temp_id     |                                      |
| title       |                                      |
| question_type|                                     |
| options     |                                      |
| answer      |                                      |
| explanation |                                      |
| keywords    |                                      |
| language    |                                      |
| ai_model    |                                      |
| user_id     |<------------------------------------+
| created_at  |
| deleted_at  |
+-------------+
```