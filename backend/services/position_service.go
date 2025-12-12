package services

import (
	"sayhi/backend/database"
)

// PositionService 位置值服务（使用数据库存储）
type PositionService struct {
	// 使用数据库存储，不再使用内存缓存
}

// NewPositionService 创建位置值服务
func NewPositionService() *PositionService {
	return &PositionService{}
}

// GetAllPositions 获取所有位置值
func (ps *PositionService) GetAllPositions() map[string][]string {
	result := make(map[string][]string)

	rows, err := database.DB.Query("SELECT position, value FROM position_values ORDER BY position, sort_order")
	if err != nil {
		return result
	}
	defer rows.Close()

	for rows.Next() {
		var position, value string
		if err := rows.Scan(&position, &value); err != nil {
			continue
		}
		result[position] = append(result[position], value)
	}

	return result
}

// GetPositionValues 获取指定位置的值
func (ps *PositionService) GetPositionValues(position string) []string {
	var values []string

	rows, err := database.DB.Query("SELECT value FROM position_values WHERE position = ? ORDER BY sort_order", position)
	if err != nil {
		return values
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			continue
		}
		values = append(values, value)
	}

	return values
}

// AddPositionValue 添加位置值
func (ps *PositionService) AddPositionValue(position string, value string) {
	// 检查是否已存在
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM position_values WHERE position = ? AND value = ?", position, value).Scan(&count)
	if count > 0 {
		return // 已存在，不重复添加
	}

	// 获取当前最大排序值
	var maxSort int
	database.DB.QueryRow("SELECT COALESCE(MAX(sort_order), 0) FROM position_values WHERE position = ?", position).Scan(&maxSort)

	// 插入新值
	database.DB.Exec("INSERT INTO position_values (position, value, sort_order) VALUES (?, ?, ?)", position, value, maxSort+1)
}

// SetPositionValues 设置位置的所有值
func (ps *PositionService) SetPositionValues(position string, values []string) {
	// 开启事务
	tx, err := database.DB.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	// 删除该位置的所有旧值
	_, err = tx.Exec("DELETE FROM position_values WHERE position = ?", position)
	if err != nil {
		return
	}

	// 插入新值
	for i, value := range values {
		_, err = tx.Exec("INSERT INTO position_values (position, value, sort_order) VALUES (?, ?, ?)", position, value, i+1)
		if err != nil {
			return
		}
	}

	// 提交事务
	tx.Commit()
}

// DeletePositionValue 删除位置值
func (ps *PositionService) DeletePositionValue(position string, value string) {
	database.DB.Exec("DELETE FROM position_values WHERE position = ? AND value = ?", position, value)
}

// UpdatePositionValue 更新位置值
func (ps *PositionService) UpdatePositionValue(position string, oldValue string, newValue string) {
	database.DB.Exec("UPDATE position_values SET value = ? WHERE position = ? AND value = ?", newValue, position, oldValue)
}
