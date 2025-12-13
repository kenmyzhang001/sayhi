package handlers

import (
	"fmt"
	"net/http"
	"sayhi/backend/models"
	"sayhi/backend/services"

	"github.com/gin-gonic/gin"
)

// TemplateHandler 模板处理器
type TemplateHandler struct {
	generator *services.TemplateGenerator
}

// NewTemplateHandler 创建模板处理器
func NewTemplateHandler(speechService *services.SpeechService) *TemplateHandler {
	return &TemplateHandler{
		generator: services.NewTemplateGenerator(speechService),
	}
}

// Generate 生成短信内容
func (h *TemplateHandler) Generate(c *gin.Context) {
	var req models.TemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	// 验证编码类型（兼容旧版本）
	if req.Encoding != "" && !isValidEncoding(req.Encoding) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的编码类型",
		})
		return
	}

	// 验证每个位置的编码类型
	if req.Encodings == nil || len(req.Encodings) == 0 {
		// 如果没有提供 Encodings，尝试使用旧的 Encoding
		if req.Encoding == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "必须提供编码配置（encodings）",
			})
			return
		}
		// 兼容旧版本：为所有选中的位置设置相同的编码
		req.Encodings = make(map[string]models.EncodingType)
		for _, pos := range req.SelectedPositions {
			req.Encodings[pos] = req.Encoding
		}
	}

	// 验证每个位置的编码是否有效
	for pos, encoding := range req.Encodings {
		if !isValidEncoding(encoding) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("位置 %s 的编码类型无效", pos),
			})
			return
		}
	}

	// 确保每个选中的位置都有编码
	for _, pos := range req.SelectedPositions {
		if _, exists := req.Encodings[pos]; !exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("位置 %s 缺少编码配置", pos),
			})
			return
		}
	}

	// 验证生成方式
	if !isValidGenerateMode(req.GenerateMode) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的生成方式",
		})
		return
	}

	response, err := h.generator.Generate(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "生成失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func isValidEncoding(encoding models.EncodingType) bool {
	switch encoding {
	case models.EncodingASCII, models.EncodingZawgyi, models.EncodingUnicode, models.EncodingOther:
		return true
	default:
		return false
	}
}

func isValidGenerateMode(mode models.GenerateMode) bool {
	switch mode {
	case models.GenerateSequential, models.GenerateRandom:
		return true
	default:
		return false
	}
}
