package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sayhi/backend/database"
	"sayhi/backend/models"
	"sayhi/backend/utils"
)

// AuthService 认证服务
type AuthService struct {
	// 使用数据库存储，不再使用内存缓存
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	return &AuthService{}
}

// hashPassword 密码加密（使用MD5，生产环境建议使用bcrypt）
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Register 注册用户
func (as *AuthService) Register(username, password string) error {
	// 检查用户名是否已存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE username = ?", username).Scan(&count)
	if err != nil {
		return errors.New("查询用户失败: " + err.Error())
	}
	if count > 0 {
		return errors.New("用户名已存在")
	}

	// 插入新用户
	hashedPassword := hashPassword(password)
	result, err := database.DB.Exec("INSERT INTO users (username, password) VALUES (?, ?)", username, hashedPassword)
	if err != nil {
		return errors.New("注册用户失败: " + err.Error())
	}

	// 验证插入成功
	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return errors.New("注册用户失败")
	}

	return nil
}

// Login 用户登录
func (as *AuthService) Login(username, password string) (string, error) {
	// 从数据库查询用户
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return "", errors.New("用户名或密码错误")
	}

	hashedPassword := hashPassword(password)
	if user.Password != hashedPassword {
		return "", errors.New("用户名或密码错误")
	}

	// 生成JWT token
	token, err := utils.GenerateToken(username, user.ID)
	if err != nil {
		return "", errors.New("生成token失败")
	}

	return token, nil
}

// ValidateToken 验证token
func (as *AuthService) ValidateToken(token string) (*models.User, error) {
	claims, err := utils.ParseToken(token)
	if err != nil {
		return nil, errors.New("无效的token")
	}

	// 从数据库查询用户
	var user models.User
	err = database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", claims.Username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &user, nil
}

// GetUser 获取用户信息
func (as *AuthService) GetUser(username string) (*models.User, error) {
	// 从数据库查询用户
	var user models.User
	err := database.DB.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).
		Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// 不返回密码
	user.Password = ""
	return &user, nil
}
