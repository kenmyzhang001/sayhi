package services

import (
	"errors"
	"sayhi/backend/database"
	"sayhi/backend/models"
)

// SpeechService 话术服务（使用数据库存储）
type SpeechService struct {
	// 使用数据库存储，不再使用内存缓存
}

// NewSpeechService 创建话术服务
func NewSpeechService() *SpeechService {
	return &SpeechService{}
}

// CreateGroup 创建话术组
func (ss *SpeechService) CreateGroup(req *models.SpeechGroupRequest) (*models.SpeechGroup, error) {
	// 检查名称是否重复
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM speech_groups WHERE name = ?", req.Name).Scan(&count)
	if err != nil {
		return nil, errors.New("查询话术组失败: " + err.Error())
	}
	if count > 0 {
		return nil, errors.New("话术组名称已存在")
	}

	// 开启事务
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, errors.New("开启事务失败: " + err.Error())
	}
	defer tx.Rollback()

	// 插入话术组
	result, err := tx.Exec("INSERT INTO speech_groups (name, description) VALUES (?, ?)", req.Name, req.Description)
	if err != nil {
		return nil, errors.New("创建话术组失败: " + err.Error())
	}

	groupID, err := result.LastInsertId()
	if err != nil {
		return nil, errors.New("获取话术组ID失败: " + err.Error())
	}

	// 插入话术内容
	for i, speech := range req.Speeches {
		_, err = tx.Exec("INSERT INTO speeches (group_id, content, sort_order) VALUES (?, ?, ?)", groupID, speech, i+1)
		if err != nil {
			return nil, errors.New("插入话术失败: " + err.Error())
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, errors.New("提交事务失败: " + err.Error())
	}

	// 返回创建的话术组
	group := &models.SpeechGroup{
		ID:          groupID,
		Name:        req.Name,
		Description: req.Description,
		Speeches:    make([]string, len(req.Speeches)),
	}
	copy(group.Speeches, req.Speeches)

	return group, nil
}

// GetGroup 获取话术组
func (ss *SpeechService) GetGroup(id int64) (*models.SpeechGroup, error) {
	var group models.SpeechGroup
	err := database.DB.QueryRow("SELECT id, name, description FROM speech_groups WHERE id = ?", id).
		Scan(&group.ID, &group.Name, &group.Description)
	if err != nil {
		return nil, errors.New("话术组不存在")
	}

	// 获取话术内容
	rows, err := database.DB.Query("SELECT content FROM speeches WHERE group_id = ? ORDER BY sort_order", id)
	if err != nil {
		return nil, errors.New("查询话术失败: " + err.Error())
	}
	defer rows.Close()

	var speeches []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			continue
		}
		speeches = append(speeches, content)
	}
	group.Speeches = speeches

	return &group, nil
}

// GetGroupByName 根据名称获取话术组
func (ss *SpeechService) GetGroupByName(name string) (*models.SpeechGroup, error) {
	var group models.SpeechGroup
	err := database.DB.QueryRow("SELECT id, name, description FROM speech_groups WHERE name = ?", name).
		Scan(&group.ID, &group.Name, &group.Description)
	if err != nil {
		return nil, errors.New("话术组不存在")
	}

	// 获取话术内容
	rows, err := database.DB.Query("SELECT content FROM speeches WHERE group_id = ? ORDER BY sort_order", group.ID)
	if err != nil {
		return nil, errors.New("查询话术失败: " + err.Error())
	}
	defer rows.Close()

	var speeches []string
	for rows.Next() {
		var content string
		if err := rows.Scan(&content); err != nil {
			continue
		}
		speeches = append(speeches, content)
	}
	group.Speeches = speeches

	return &group, nil
}

// GetAllGroups 获取所有话术组
func (ss *SpeechService) GetAllGroups() []models.SpeechGroup {
	rows, err := database.DB.Query("SELECT id, name, description FROM speech_groups ORDER BY id")
	if err != nil {
		return []models.SpeechGroup{}
	}
	defer rows.Close()

	var groups []models.SpeechGroup
	for rows.Next() {
		var group models.SpeechGroup
		if err := rows.Scan(&group.ID, &group.Name, &group.Description); err != nil {
			continue
		}

		// 获取话术内容
		speechRows, err := database.DB.Query("SELECT content FROM speeches WHERE group_id = ? ORDER BY sort_order", group.ID)
		if err == nil {
			var speeches []string
			for speechRows.Next() {
				var content string
				if err := speechRows.Scan(&content); err != nil {
					continue
				}
				speeches = append(speeches, content)
			}
			group.Speeches = speeches
			speechRows.Close()
		}

		groups = append(groups, group)
	}

	return groups
}

// UpdateGroup 更新话术组
func (ss *SpeechService) UpdateGroup(id int64, req *models.SpeechGroupUpdateRequest) (*models.SpeechGroup, error) {
	// 检查话术组是否存在
	var currentName string
	err := database.DB.QueryRow("SELECT name FROM speech_groups WHERE id = ?", id).Scan(&currentName)
	if err != nil {
		return nil, errors.New("话术组不存在")
	}

	// 如果更新名称，检查是否与其他组重复
	if req.Name != "" && req.Name != currentName {
		var count int
		err := database.DB.QueryRow("SELECT COUNT(*) FROM speech_groups WHERE name = ? AND id != ?", req.Name, id).Scan(&count)
		if err != nil {
			return nil, errors.New("查询话术组失败: " + err.Error())
		}
		if count > 0 {
			return nil, errors.New("话术组名称已存在")
		}
	}

	// 开启事务
	tx, err := database.DB.Begin()
	if err != nil {
		return nil, errors.New("开启事务失败: " + err.Error())
	}
	defer tx.Rollback()

	// 更新话术组基本信息
	if req.Name != "" || req.Description != "" {
		if req.Name != "" && req.Description != "" {
			_, err = tx.Exec("UPDATE speech_groups SET name = ?, description = ? WHERE id = ?", req.Name, req.Description, id)
		} else if req.Name != "" {
			_, err = tx.Exec("UPDATE speech_groups SET name = ? WHERE id = ?", req.Name, id)
		} else {
			_, err = tx.Exec("UPDATE speech_groups SET description = ? WHERE id = ?", req.Description, id)
		}
		if err != nil {
			return nil, errors.New("更新话术组失败: " + err.Error())
		}
	}

	// 如果更新话术内容
	if len(req.Speeches) > 0 {
		// 删除旧的话术
		_, err = tx.Exec("DELETE FROM speeches WHERE group_id = ?", id)
		if err != nil {
			return nil, errors.New("删除旧话术失败: " + err.Error())
		}

		// 插入新话术
		for i, speech := range req.Speeches {
			_, err = tx.Exec("INSERT INTO speeches (group_id, content, sort_order) VALUES (?, ?, ?)", id, speech, i+1)
			if err != nil {
				return nil, errors.New("插入话术失败: " + err.Error())
			}
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return nil, errors.New("提交事务失败: " + err.Error())
	}

	// 返回更新后的话术组
	return ss.GetGroup(id)
}

// DeleteGroup 删除话术组
func (ss *SpeechService) DeleteGroup(id int64) error {
	// 检查话术组是否存在
	var count int
	err := database.DB.QueryRow("SELECT COUNT(*) FROM speech_groups WHERE id = ?", id).Scan(&count)
	if err != nil {
		return errors.New("查询话术组失败: " + err.Error())
	}
	if count == 0 {
		return errors.New("话术组不存在")
	}

	// 删除话术组（由于外键约束，会自动删除关联的话术）
	_, err = database.DB.Exec("DELETE FROM speech_groups WHERE id = ?", id)
	if err != nil {
		return errors.New("删除话术组失败: " + err.Error())
	}

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
