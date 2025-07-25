package services

import (
	"CodeQuizAI/dao"
	"CodeQuizAI/models"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ListUsers 管理员获取用户列表（业务逻辑）
func ListUsers(
	c *gin.Context,
	page, pageSize int,
	username string,
) (userList []map[string]interface{}, pagination map[string]int, err error) {
	// 1. 校验分页参数（仅保留参数合法性检查）
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	// 2. 调用DAO查询数据
	query := dao.Q.User.WithContext(c).Offset(offset).Limit(pageSize)
	if username != "" {
		query = query.Where(dao.User.Username.Like("%" + username + "%"))
	}

	// 3. 获取用户列表和总条数
	var users []*models.User
	users, err = query.Find()
	if err != nil {
		return nil, nil, errors.New("查询用户列表失败")
	}

	total, err := dao.Q.User.WithContext(c).Count()
	if err != nil {
		return nil, nil, errors.New("查询用户总数失败")
	}

	// 4. 数据脱敏（隐藏敏感字段）
	userList = make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		userList = append(userList, map[string]interface{}{
			"id":         u.ID,
			"username":   u.Username,
			"role":       u.Role,
			"created_at": u.CreatedAt,
			"updated_at": u.UpdatedAt,
		})
	}

	// 5. 构建分页信息
	pagination = map[string]int{
		"total":      int(total),
		"page":       page,
		"page_size":  pageSize,
		"total_page": (int(total) + pageSize - 1) / pageSize,
	}

	return userList, pagination, nil
}

// UpdateUser 允许用户更新自己的用户名和密码
func UpdateUser(c *gin.Context, userID int64, newUsername, newPassword string) error {
	// 1. 校验用户是否存在（通过ID查询）
	_, err := dao.Q.User.WithContext(c).Where(dao.User.ID.Eq(userID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("用户不存在")
		}
		return errors.New("查询用户失败: " + err.Error())
	}

	// 2. 构建更新字段（仅处理非空参数）
	updates := make(map[string]interface{})

	// 处理用户名更新（需校验唯一性）
	if newUsername != "" {
		// 排除当前用户ID，检查用户名是否已被其他用户占用
		exist, _ := dao.Q.User.WithContext(c).
			Where(dao.User.Username.Eq(newUsername), dao.User.ID.Neq(userID)).
			First()
		if exist != nil {
			return errors.New("用户名已被占用")
		}
		updates["username"] = newUsername
	}

	// 处理密码更新（需加密）
	if newPassword != "" {
		// 计算 SHA-256 哈希
		hash := sha256.Sum256([]byte(newPassword))
		// 转为十六进制字符串
		hashStr := hex.EncodeToString(hash[:])
		updates["password_hash"] = hashStr
	}

	// 3. 执行更新操作
	if len(updates) > 0 {
		_, err = dao.Q.User.WithContext(c).
			Where(dao.User.ID.Eq(userID)).
			Updates(updates)
		if err != nil {
			return errors.New("更新用户信息失败: " + err.Error())
		}
	}

	return nil
}
