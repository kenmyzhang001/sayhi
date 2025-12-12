package services

import (
	"sync"
)

// PositionService 位置值服务（内存存储，生产环境应使用数据库）
type PositionService struct {
	mu        sync.RWMutex
	positions map[string][]string // position -> values
}

// NewPositionService 创建位置值服务
func NewPositionService() *PositionService {
	return &PositionService{
		positions: make(map[string][]string),
	}
}

// GetAllPositions 获取所有位置值
func (ps *PositionService) GetAllPositions() map[string][]string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	result := make(map[string][]string)
	for k, v := range ps.positions {
		result[k] = make([]string, len(v))
		copy(result[k], v)
	}
	return result
}

// GetPositionValues 获取指定位置的值
func (ps *PositionService) GetPositionValues(position string) []string {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	values, exists := ps.positions[position]
	if !exists {
		return []string{}
	}

	result := make([]string, len(values))
	copy(result, values)
	return result
}

// AddPositionValue 添加位置值
func (ps *PositionService) AddPositionValue(position string, value string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if ps.positions[position] == nil {
		ps.positions[position] = []string{}
	}

	// 检查是否已存在
	for _, v := range ps.positions[position] {
		if v == value {
			return // 已存在，不重复添加
		}
	}

	ps.positions[position] = append(ps.positions[position], value)
}

// SetPositionValues 设置位置的所有值
func (ps *PositionService) SetPositionValues(position string, values []string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ps.positions[position] = make([]string, len(values))
	copy(ps.positions[position], values)
}

// DeletePositionValue 删除位置值
func (ps *PositionService) DeletePositionValue(position string, value string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	values, exists := ps.positions[position]
	if !exists {
		return
	}

	var newValues []string
	for _, v := range values {
		if v != value {
			newValues = append(newValues, v)
		}
	}

	ps.positions[position] = newValues
}

// UpdatePositionValue 更新位置值
func (ps *PositionService) UpdatePositionValue(position string, oldValue string, newValue string) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	values, exists := ps.positions[position]
	if !exists {
		return
	}

	for i, v := range values {
		if v == oldValue {
			values[i] = newValue
			break
		}
	}

	ps.positions[position] = values
}
