package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SpeechResolver 话术解析器接口
type SpeechResolver interface {
	ResolveSpeechGroup(nameOrID string) ([]string, error)
}

// ParseTemplate 解析模板，返回位置列表和原始模板结构
func ParseTemplate(template string) ([]string, []string, error) {
	// 匹配所有括号内容
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(template, -1)
	
	if len(matches) == 0 {
		return nil, nil, fmt.Errorf("模板格式错误：未找到括号位置")
	}

	var positions []string
	var rawTemplate []string
	
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}
		content := match[1]
		positions = append(positions, content)
		rawTemplate = append(rawTemplate, match[0]) // 保存原始括号格式
	}

	return positions, rawTemplate, nil
}

// ExpandRange 展开范围值，如 "3-10" 返回 ["3", "4", "5", ..., "10"]
func ExpandRange(rangeStr string) ([]string, error) {
	parts := strings.Split(rangeStr, "-")
	if len(parts) != 2 {
		return []string{rangeStr}, nil // 不是范围，返回原值
	}

	start, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

	if err1 != nil || err2 != nil {
		return []string{rangeStr}, nil // 解析失败，返回原值
	}

	if start > end {
		return nil, fmt.Errorf("无效的范围：起始值 %d 大于结束值 %d", start, end)
	}

	var result []string
	for i := start; i <= end; i++ {
		result = append(result, strconv.Itoa(i))
	}

	return result, nil
}

// ResolvePositionValues 解析位置值，支持固定值和范围值
func ResolvePositionValues(positionStr string, configValues []string) ([]string, error) {
	// 先尝试解析为范围
	rangeValues, err := ExpandRange(positionStr)
	if err != nil {
		return nil, err
	}

	// 如果是范围值（长度>1），返回范围值
	if len(rangeValues) > 1 {
		return rangeValues, nil
	}

	// 如果是固定值，检查是否在配置中
	if len(configValues) > 0 {
		return configValues, nil
	}

	// 如果配置为空，使用模板中的固定值
	return []string{positionStr}, nil
}

