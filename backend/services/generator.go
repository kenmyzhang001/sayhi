package services

import (
	"fmt"
	"math/rand"
	"sayhi/backend/models"
	"sayhi/backend/utils"
	"strings"
	"time"
)

const MaxCharsPerSMS = 70

// TemplateGenerator 模板生成器
type TemplateGenerator struct {
	speechService *SpeechService
}

// NewTemplateGenerator 创建新的模板生成器
func NewTemplateGenerator(speechService *SpeechService) *TemplateGenerator {
	return &TemplateGenerator{
		speechService: speechService,
	}
}

// Generate 生成短信内容
func (tg *TemplateGenerator) Generate(req *models.TemplateRequest) (*models.GenerateResponse, error) {
	// 解析模板
	positions, rawTemplate, err := utils.ParseTemplate(req.Template)
	if err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, fmt.Errorf("模板为空")
	}

	// 解析位置值（支持话术组）
	var speechGroups map[string]string
	if req.SpeechGroups != nil {
		speechGroups = req.SpeechGroups
	}
	positionValues, err := tg.resolvePositionValues(positions, req.Positions, speechGroups)
	if err != nil {
		return nil, err
	}

	// 生成所有组合
	var results []models.GeneratedResult

	if req.GenerateMode == models.GenerateSequential {
		results = tg.generateSequential(rawTemplate, positionValues, req.Encoding)
	} else {
		results = tg.generateRandom(rawTemplate, positionValues, req.Encoding)
	}

	// 统计超出数量
	exceededCount := 0
	for _, result := range results {
		if result.IsExceeded {
			exceededCount++
		}
	}

	return &models.GenerateResponse{
		Results:      results,
		TotalCount:   len(results),
		ExceededCount: exceededCount,
	}, nil
}

// resolvePositionValues 解析所有位置的值（支持话术组）
func (tg *TemplateGenerator) resolvePositionValues(positions []string, config models.PositionConfig, speechGroups map[string]string) ([][]string, error) {
	var positionValues [][]string

	// 位置标识映射
	positionKeys := []string{"a", "b", "c", "d"}

	for i, pos := range positions {
		var configValues []string
		positionKey := ""
		if i < len(positionKeys) {
			positionKey = positionKeys[i]
		}

		// 优先检查是否指定了话术组
		if speechGroups != nil && positionKey != "" {
			if speechGroupName, exists := speechGroups[positionKey]; exists {
				speeches, err := tg.speechService.GetGroupSpeeches(speechGroupName)
				if err == nil {
					positionValues = append(positionValues, speeches)
					continue
				}
			}
		}

		// 使用配置的位置值
		switch i {
		case 0: // a
			configValues = config.A
		case 1: // b
			configValues = config.B
		case 2: // c
			configValues = config.C
		case 3: // d
			configValues = config.D
		default:
			// 支持更多位置
			configValues = []string{pos}
		}

		values, err := utils.ResolvePositionValues(pos, configValues)
		if err != nil {
			return nil, fmt.Errorf("解析位置 %d 失败: %v", i, err)
		}

		positionValues = append(positionValues, values)
	}

	return positionValues, nil
}

// generateSequential 顺序生成
func (tg *TemplateGenerator) generateSequential(rawTemplate []string, positionValues [][]string, encoding models.EncodingType) []models.GeneratedResult {
	var results []models.GeneratedResult
	combinations := tg.generateCombinations(positionValues)

	for _, combo := range combinations {
		content := tg.buildContent(rawTemplate, combo)
		charCount := utils.CountChars(content, encoding)
		isExceeded := utils.IsExceeded(charCount, MaxCharsPerSMS)
		exceededChars := 0
		if isExceeded {
			exceededChars = charCount - MaxCharsPerSMS
		}

		results = append(results, models.GeneratedResult{
			Content:      content,
			CharCount:    charCount,
			IsExceeded:   isExceeded,
			ExceededChars: exceededChars,
		})
	}

	return results
}

// generateRandom 随机生成
func (tg *TemplateGenerator) generateRandom(rawTemplate []string, positionValues [][]string, encoding models.EncodingType) []models.GeneratedResult {
	// 先生成所有组合
	combinations := tg.generateCombinations(positionValues)
	
	// 随机打乱
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(combinations), func(i, j int) {
		combinations[i], combinations[j] = combinations[j], combinations[i]
	})

	var results []models.GeneratedResult
	for _, combo := range combinations {
		// 随机打乱位置顺序
		shuffledCombo := make([]string, len(combo))
		copy(shuffledCombo, combo)
		rand.Shuffle(len(shuffledCombo), func(i, j int) {
			shuffledCombo[i], shuffledCombo[j] = shuffledCombo[j], shuffledCombo[i]
		})

		content := tg.buildContentRandom(rawTemplate, shuffledCombo)
		charCount := utils.CountChars(content, encoding)
		isExceeded := utils.IsExceeded(charCount, MaxCharsPerSMS)
		exceededChars := 0
		if isExceeded {
			exceededChars = charCount - MaxCharsPerSMS
		}

		results = append(results, models.GeneratedResult{
			Content:      content,
			CharCount:    charCount,
			IsExceeded:   isExceeded,
			ExceededChars: exceededChars,
		})
	}

	return results
}

// generateCombinations 生成所有组合（笛卡尔积）
func (tg *TemplateGenerator) generateCombinations(positionValues [][]string) [][]string {
	if len(positionValues) == 0 {
		return [][]string{}
	}

	if len(positionValues) == 1 {
		var result [][]string
		for _, val := range positionValues[0] {
			result = append(result, []string{val})
		}
		return result
	}

	// 递归生成组合
	rest := tg.generateCombinations(positionValues[1:])
	var result [][]string

	for _, val := range positionValues[0] {
		for _, combo := range rest {
			newCombo := append([]string{val}, combo...)
			result = append(result, newCombo)
		}
	}

	return result
}

// buildContent 构建内容（顺序）
func (tg *TemplateGenerator) buildContent(rawTemplate []string, values []string) string {
	content := strings.Join(rawTemplate, "")
	for i, val := range values {
		if i < len(rawTemplate) {
			// 替换括号内容
			old := rawTemplate[i]
			content = strings.Replace(content, old, val, 1)
		}
	}
	// 移除剩余的括号
	content = strings.ReplaceAll(content, "(", "")
	content = strings.ReplaceAll(content, ")", "")
	return content
}

// buildContentRandom 构建内容（随机顺序）
func (tg *TemplateGenerator) buildContentRandom(rawTemplate []string, values []string) string {
	// 随机顺序：直接使用值，不按模板顺序
	return strings.Join(values, " ")
}

