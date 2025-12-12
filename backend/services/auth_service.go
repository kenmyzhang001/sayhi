package services

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"sayhi/backend/models"
	"sayhi/backend/utils"
	"sync"
)

// AuthService 认证服务
type AuthService struct {
	mu    sync.RWMutex
	users map[string]*models.User // username -> User
}

// NewAuthService 创建认证服务
func NewAuthService() *AuthService {
	service := &AuthService{
		users: make(map[string]*models.User),
	}

	// 初始化默认管理员账号
	service.Register("admin", "admin123")
	service.Register("user", "user123")

	return service
}

// hashPassword 密码加密（使用MD5，生产环境建议使用bcrypt）
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// Register 注册用户
func (as *AuthService) Register(username, password string) error {
	as.mu.Lock()
	defer as.mu.Unlock()

	if _, exists := as.users[username]; exists {
		return errors.New("用户名已存在")
	}

	as.users[username] = &models.User{
		ID:       int64(len(as.users) + 1),
		Username: username,
		Password: hashPassword(password),
	}

	return nil
}

// Login 用户登录
func (as *AuthService) Login(username, password string) (string, error) {
	as.mu.RLock()
	user, exists := as.users[username]
	as.mu.RUnlock()

	if !exists {
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

	as.mu.RLock()
	user, exists := as.users[claims.Username]
	as.mu.RUnlock()

	if !exists {
		return nil, errors.New("用户不存在")
	}

	return user, nil
}

// GetUser 获取用户信息
func (as *AuthService) GetUser(username string) (*models.User, error) {
	as.mu.RLock()
	defer as.mu.RUnlock()

	user, exists := as.users[username]
	if !exists {
		return nil, errors.New("用户不存在")
	}

	// 不返回密码
	user.Password = ""
	return user, nil
}

