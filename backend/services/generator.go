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
	var positionKeys []string
	var positionValues [][]string

	// 如果提供了模板，解析模板
	if req.Template != "" {
		positions, _, err := utils.ParseTemplate(req.Template)
		if err != nil {
			return nil, err
		}
		if len(positions) == 0 {
			return nil, fmt.Errorf("模板为空")
		}

		// 从模板解析位置
		positionKeys = []string{"a", "b", "c", "d"}
		if len(positions) > len(positionKeys) {
			// 扩展位置键
			for i := len(positionKeys); i < len(positions); i++ {
				positionKeys = append(positionKeys, string(rune('a'+i)))
			}
		}
		positionKeys = positionKeys[:len(positions)]

		// 解析位置值（支持话术组）
		var speechGroups map[string]string
		if req.SpeechGroups != nil {
			speechGroups = req.SpeechGroups
		}
		positionValues, err = tg.resolvePositionValues(positions, req.Positions, speechGroups)
		if err != nil {
			return nil, err
		}
	} else {
		// 如果没有模板，使用选择的位置
		if len(req.SelectedPositions) == 0 {
			return nil, fmt.Errorf("请至少选择一个位置")
		}

		positionKeys = req.SelectedPositions
		var speechGroups map[string]string
		if req.SpeechGroups != nil {
			speechGroups = req.SpeechGroups
		}

		// 根据选择的位置获取值
		positionValues = make([][]string, 0, len(positionKeys))
		for _, posKey := range positionKeys {
			var values []string

			// 优先检查是否指定了话术组
			if speechGroups != nil {
				if speechGroupName, exists := speechGroups[posKey]; exists {
					speeches, err := tg.speechService.GetGroupSpeeches(speechGroupName)
					if err == nil {
						positionValues = append(positionValues, speeches)
						continue
					}
				}
			}

			// 使用配置的位置值
			switch posKey {
			case "a":
				values = req.Positions.A
			case "b":
				values = req.Positions.B
			case "c":
				values = req.Positions.C
			case "d":
				values = req.Positions.D
			case "e":
				values = req.Positions.E
			case "f":
				values = req.Positions.F
			case "g":
				values = req.Positions.G
			case "h":
				values = req.Positions.H
			case "i":
				values = req.Positions.I
			case "j":
				values = req.Positions.J
			default:
				return nil, fmt.Errorf("不支持的位置: %s", posKey)
			}

			if len(values) == 0 {
				return nil, fmt.Errorf("位置 %s 没有配置值或话术组", posKey)
			}

			positionValues = append(positionValues, values)
		}
	}

	// 生成所有组合
	var results []models.GeneratedResult

	if req.GenerateMode == models.GenerateSequential {
		results = tg.generateSequential(positionKeys, positionValues, req.Encodings)
	} else {
		results = tg.generateRandom(positionKeys, positionValues, req.Encodings)
	}

	// 统计超出数量
	exceededCount := 0
	for _, result := range results {
		if result.IsExceeded {
			exceededCount++
		}
	}

	return &models.GenerateResponse{
		Results:       results,
		TotalCount:    len(results),
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
func (tg *TemplateGenerator) generateSequential(positionKeys []string, positionValues [][]string, encodings map[string]models.EncodingType) []models.GeneratedResult {
	var results []models.GeneratedResult
	combinations := tg.generateCombinations(positionValues)

	for _, combo := range combinations {
		content := tg.buildContentFromValues(positionKeys, combo)
		// 使用每个位置对应的编码计算字符数
		charCount := tg.countCharsWithPositionEncodings(positionKeys, combo, encodings)
		isExceeded := utils.IsExceeded(charCount, MaxCharsPerSMS)
		exceededChars := 0
		if isExceeded {
			exceededChars = charCount - MaxCharsPerSMS
		}

		results = append(results, models.GeneratedResult{
			Content:       content,
			CharCount:     charCount,
			IsExceeded:    isExceeded,
			ExceededChars: exceededChars,
		})
	}

	return results
}

// generateRandom 随机生成
func (tg *TemplateGenerator) generateRandom(positionKeys []string, positionValues [][]string, encodings map[string]models.EncodingType) []models.GeneratedResult {
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
		shuffledKeys := make([]string, len(positionKeys))
		shuffledCombo := make([]string, len(combo))
		shuffledEncodings := make(map[string]models.EncodingType)
		copy(shuffledKeys, positionKeys)
		copy(shuffledCombo, combo)

		// 同时打乱键和值，保持对应关系
		rand.Shuffle(len(shuffledKeys), func(i, j int) {
			shuffledKeys[i], shuffledKeys[j] = shuffledKeys[j], shuffledKeys[i]
			shuffledCombo[i], shuffledCombo[j] = shuffledCombo[j], shuffledCombo[i]
		})

		// 更新编码映射以匹配打乱后的位置
		for i, key := range shuffledKeys {
			if encoding, exists := encodings[positionKeys[i]]; exists {
				shuffledEncodings[key] = encoding
			}
		}

		content := tg.buildContentFromValues(shuffledKeys, shuffledCombo)
		// 使用每个位置对应的编码计算字符数
		charCount := tg.countCharsWithPositionEncodings(shuffledKeys, shuffledCombo, shuffledEncodings)
		isExceeded := utils.IsExceeded(charCount, MaxCharsPerSMS)
		exceededChars := 0
		if isExceeded {
			exceededChars = charCount - MaxCharsPerSMS
		}

		results = append(results, models.GeneratedResult{
			Content:       content,
			CharCount:     charCount,
			IsExceeded:    isExceeded,
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

// buildContentFromValues 根据位置键和值构建内容
func (tg *TemplateGenerator) buildContentFromValues(positionKeys []string, values []string) string {
	// 按位置顺序组合：a b c d 的值
	return strings.Join(values, " ")
}

// countCharsWithPositionEncodings 使用每个位置对应的编码计算总字符数
// 每个位置的值使用该位置对应的编码来计算字符数，然后相加
// 空格按ASCII编码计算（每个空格1个字符）
func (tg *TemplateGenerator) countCharsWithPositionEncodings(positionKeys []string, values []string, encodings map[string]models.EncodingType) int {
	totalCount := 0
	spaceCount := len(values) - 1 // 空格的数量

	for i, key := range positionKeys {
		if i < len(values) {
			encoding := models.EncodingUnicode // 默认编码
			if enc, exists := encodings[key]; exists {
				encoding = enc
			}
			// 计算该位置值的字符数
			totalCount += utils.CountChars(values[i], encoding)
		}
	}

	// 添加空格（空格按ASCII计算，每个空格1个字符）
	totalCount += spaceCount

	return totalCount
}
