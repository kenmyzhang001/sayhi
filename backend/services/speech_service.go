package services

import (
	"errors"
	"sayhi/backend/models"
	"sync"
)

// SpeechService 话术服务（内存存储，生产环境应使用数据库）
type SpeechService struct {
	mu     sync.RWMutex
	groups map[int64]*models.SpeechGroup
	nextID int64
}

// NewSpeechService 创建话术服务
func NewSpeechService() *SpeechService {
	service := &SpeechService{
		groups: make(map[int64]*models.SpeechGroup),
		nextID: 1,
	}

	// 初始化示例话术组
	service.initDefaultGroups()

	return service
}

// initDefaultGroups 初始化默认话术组
func (ss *SpeechService) initDefaultGroups() {
	// 示例：数字范围话术组
	ss.CreateGroup(&models.SpeechGroupRequest{
		Name:        "数字3-10",
		Description: "数字范围3到10",
		Speeches:    []string{"3", "4", "5", "6", "7", "8", "9", "10"},
	})

	// 示例：问候语话术组
	ss.CreateGroup(&models.SpeechGroupRequest{
		Name:        "问候语",
		Description: "常用问候语",
		Speeches:    []string{"您好", "早上好", "下午好", "晚上好", "欢迎"},
	})
}

// CreateGroup 创建话术组
func (ss *SpeechService) CreateGroup(req *models.SpeechGroupRequest) (*models.SpeechGroup, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	// 检查名称是否重复
	for _, group := range ss.groups {
		if group.Name == req.Name {
			return nil, errors.New("话术组名称已存在")
		}
	}

	group := &models.SpeechGroup{
		ID:          ss.nextID,
		Name:        req.Name,
		Description: req.Description,
		Speeches:    make([]string, len(req.Speeches)),
	}

	copy(group.Speeches, req.Speeches)

	ss.groups[ss.nextID] = group
	ss.nextID++

	return group, nil
}

// GetGroup 获取话术组
func (ss *SpeechService) GetGroup(id int64) (*models.SpeechGroup, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	group, exists := ss.groups[id]
	if !exists {
		return nil, errors.New("话术组不存在")
	}

	// 返回副本
	result := *group
	result.Speeches = make([]string, len(group.Speeches))
	copy(result.Speeches, group.Speeches)

	return &result, nil
}

// GetGroupByName 根据名称获取话术组
func (ss *SpeechService) GetGroupByName(name string) (*models.SpeechGroup, error) {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	for _, group := range ss.groups {
		if group.Name == name {
			result := *group
			result.Speeches = make([]string, len(group.Speeches))
			copy(result.Speeches, group.Speeches)
			return &result, nil
		}
	}

	return nil, errors.New("话术组不存在")
}

// GetAllGroups 获取所有话术组
func (ss *SpeechService) GetAllGroups() []models.SpeechGroup {
	ss.mu.RLock()
	defer ss.mu.RUnlock()

	groups := make([]models.SpeechGroup, 0, len(ss.groups))
	for _, group := range ss.groups {
		result := *group
		result.Speeches = make([]string, len(group.Speeches))
		copy(result.Speeches, group.Speeches)
		groups = append(groups, result)
	}

	return groups
}

// UpdateGroup 更新话术组
func (ss *SpeechService) UpdateGroup(id int64, req *models.SpeechGroupUpdateRequest) (*models.SpeechGroup, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	group, exists := ss.groups[id]
	if !exists {
		return nil, errors.New("话术组不存在")
	}

	// 如果更新名称，检查是否与其他组重复
	if req.Name != "" && req.Name != group.Name {
		for _, g := range ss.groups {
			if g.ID != id && g.Name == req.Name {
				return nil, errors.New("话术组名称已存在")
			}
		}
		group.Name = req.Name
	}

	if req.Description != "" {
		group.Description = req.Description
	}

	if len(req.Speeches) > 0 {
		group.Speeches = make([]string, len(req.Speeches))
		copy(group.Speeches, req.Speeches)
	}

	result := *group
	result.Speeches = make([]string, len(group.Speeches))
	copy(result.Speeches, group.Speeches)

	return &result, nil
}

// DeleteGroup 删除话术组
func (ss *SpeechService) DeleteGroup(id int64) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if _, exists := ss.groups[id]; !exists {
		return errors.New("话术组不存在")
	}

	delete(ss.groups, id)
	return nil
}

// GetGroupSpeeches 获取话术组的所有话术
func (ss *SpeechService) GetGroupSpeeches(nameOrID string) ([]string, error) {
	// 先尝试按ID查找
	if id, err := parseInt64(nameOrID); err == nil {
		group, err := ss.GetGroup(id)
		if err == nil {
			return group.Speeches, nil
		}
	}

	// 再尝试按名称查找
	group, err := ss.GetGroupByName(nameOrID)
	if err != nil {
		return nil, errors.New("话术组不存在: " + nameOrID)
	}

	return group.Speeches, nil
}

// parseInt64 尝试将字符串转换为int64
func parseInt64(s string) (int64, error) {
	var result int64
	var err error
	for _, char := range s {
		if char < '0' || char > '9' {
			return 0, errors.New("not a number")
		}
		result = result*10 + int64(char-'0')
	}
	return result, err
}

