package models

// SpeechGroup 话术组
type SpeechGroup struct {
	ID          int64    `json:"id"`
	Name        string   `json:"name" binding:"required"`        // 话术组名称
	Description string   `json:"description"`                    // 描述
	Speeches    []string `json:"speeches" binding:"required"`    // 话术列表
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
}

// SpeechGroupRequest 话术组请求
type SpeechGroupRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description string   `json:"description"`
	Speeches    []string `json:"speeches" binding:"required,min=1"`
}

// SpeechGroupUpdateRequest 话术组更新请求
type SpeechGroupUpdateRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Speeches    []string `json:"speeches"`
}

// SpeechGroupListResponse 话术组列表响应
type SpeechGroupListResponse struct {
	Groups []SpeechGroup `json:"groups"`
	Total  int           `json:"total"`
}

