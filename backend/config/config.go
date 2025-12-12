package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string
	Port string
	Mode string // debug, release, test
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Type     string // mysql, postgresql, sqlite
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
	MaxIdle  int    // 最大空闲连接数
	MaxOpen  int    // 最大打开连接数
	DSN      string // 完整连接字符串（如果提供则优先使用）
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string
	ExpireTime int // 过期时间（小时）
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() *Config {
	// 加载 .env 文件（如果存在）
	// 忽略错误，因为 .env 文件是可选的
	if err := godotenv.Load(); err != nil {
		// 尝试从 backend 目录加载
		_ = godotenv.Load("../.env")
		// 如果都不存在，使用环境变量或默认值
	}

	config := &Config{
		Server: ServerConfig{
			Host: getEnv("SERVER_HOST", "0.0.0.0"),
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Type:     getEnv("DB_TYPE", "mysql"),
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			Database: getEnv("DB_NAME", "sayhi"),
			Charset:  getEnv("DB_CHARSET", "utf8mb4"),
			MaxIdle:  getEnvAsInt("DB_MAX_IDLE", 10),
			MaxOpen:  getEnvAsInt("DB_MAX_OPEN", 100),
			DSN:      getEnv("DB_DSN", ""), // 如果设置了DSN，则优先使用
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "sayhi-secret-key-change-in-production"),
			ExpireTime: getEnvAsInt("JWT_EXPIRE_TIME", 24), // 24小时
		},
	}

	AppConfig = config
	return config
}

// GetDSN 获取数据库连接字符串
func (c *DatabaseConfig) GetDSN() string {
	if c.DSN != "" {
		return c.DSN
	}

	switch c.Type {
	case "mysql":
		return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + c.Port + ")/" + c.Database + "?charset=" + c.Charset + "&parseTime=True&loc=Local"
	case "postgresql", "postgres":
		return "host=" + c.Host + " port=" + c.Port + " user=" + c.User + " password=" + c.Password + " dbname=" + c.Database + " sslmode=disable"
	case "sqlite":
		if c.Database == "" {
			return "sayhi.db"
		}
		return c.Database
	default:
		return ""
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取环境变量并转换为整数
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}
