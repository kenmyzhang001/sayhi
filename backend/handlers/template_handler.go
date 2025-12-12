package handlers

import (
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

	// 验证编码类型
	if !isValidEncoding(req.Encoding) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的编码类型",
		})
		return
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

