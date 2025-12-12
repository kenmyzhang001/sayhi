package database

import (
	"database/sql"
	"fmt"
	"sayhi/backend/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

// InitDB 初始化数据库连接
func InitDB(cfg *config.DatabaseConfig) error {
	dsn := cfg.GetDSN()
	if dsn == "" {
		return fmt.Errorf("数据库配置错误：无法生成 DSN")
	}

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("打开数据库连接失败: %w", err)
	}

	// 设置连接池参数
	DB.SetMaxIdleConns(cfg.MaxIdle)
	DB.SetMaxOpenConns(cfg.MaxOpen)
	DB.SetConnMaxLifetime(time.Hour)

	// 测试连接
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("数据库连接测试失败: %w", err)
	}

	return nil
}

// CloseDB 关闭数据库连接
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
