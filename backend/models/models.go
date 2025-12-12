package models

// EncodingType 字符编码类型
type EncodingType string

const (
	EncodingASCII  EncodingType = "ASCII"
	EncodingZawgyi EncodingType = "Zawgyi"
	EncodingUnicode EncodingType = "Unicode"
	EncodingOther   EncodingType = "Other"
)

// GenerateMode 生成方式
type GenerateMode string

const (
	GenerateSequential GenerateMode = "sequential"
	GenerateRandom     GenerateMode = "random"
)

// TemplateRequest 模板生成请求
type TemplateRequest struct {
	Template     string        `json:"template,omitempty"` // 模板（可选，如果不提供则根据位置自动生成）
	Encoding     EncodingType  `json:"encoding" binding:"required"`
	GenerateMode GenerateMode  `json:"generateMode" binding:"required"`
	Positions    PositionConfig `json:"positions" binding:"required"`
	SpeechGroups map[string]string `json:"speechGroups,omitempty"` // 位置 -> 话术组名称或ID的映射
	SelectedPositions []string `json:"selectedPositions,omitempty"` // 选择的位置（如 ["a", "b", "c", "d"]）
}

// PositionConfig 位置配置
type PositionConfig struct {
	A []string `json:"a"`
	B []string `json:"b"`
	C []string `json:"c"`
	D []string `json:"d"`
}

// GeneratedResult 生成结果
type GeneratedResult struct {
	Content      string `json:"content"`
	CharCount    int    `json:"charCount"`
	IsExceeded   bool   `json:"isExceeded"`
	ExceededChars int   `json:"exceededChars"`
}

// GenerateResponse 生成响应
type GenerateResponse struct {
	Results      []GeneratedResult `json:"results"`
	TotalCount   int               `json:"totalCount"`
	ExceededCount int              `json:"exceededCount"`
}

// PositionValue 位置值配置
type PositionValue struct {
	ID       int64  `json:"id"`
	Position string `json:"position" binding:"required"` // a, b, c, d
	Value    string `json:"value" binding:"required"`
}

// PositionValueRequest 位置值请求
type PositionValueRequest struct {
	Position string `json:"position" binding:"required"`
	Value    string `json:"value" binding:"required"`
}

// PositionValueListResponse 位置值列表响应
type PositionValueListResponse struct {
	Positions map[string][]string `json:"positions"`
}

