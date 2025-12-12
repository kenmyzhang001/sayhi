package handlers

import (
	"net/http"
	"sayhi/backend/models"
	"sayhi/backend/services"

	"github.com/gin-gonic/gin"
)

// PositionHandler 位置值处理器
type PositionHandler struct {
	service *services.PositionService
}

// NewPositionHandler 创建位置值处理器
func NewPositionHandler(service *services.PositionService) *PositionHandler {
	return &PositionHandler{
		service: service,
	}
}

// GetAllPositions 获取所有位置值
func (h *PositionHandler) GetAllPositions(c *gin.Context) {
	positions := h.service.GetAllPositions()
	c.JSON(http.StatusOK, models.PositionValueListResponse{
		Positions: positions,
	})
}

// GetPositionValues 获取指定位置的值
func (h *PositionHandler) GetPositionValues(c *gin.Context) {
	position := c.Param("position")
	if !isValidPosition(position) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的位置标识",
		})
		return
	}

	values := h.service.GetPositionValues(position)
	c.JSON(http.StatusOK, gin.H{
		"position": position,
		"values":   values,
	})
}

// AddPositionValue 添加位置值
func (h *PositionHandler) AddPositionValue(c *gin.Context) {
	var req models.PositionValueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	if !isValidPosition(req.Position) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的位置标识",
		})
		return
	}

	h.service.AddPositionValue(req.Position, req.Value)
	c.JSON(http.StatusOK, gin.H{
		"message": "添加成功",
	})
}

// SetPositionValues 设置位置的所有值
func (h *PositionHandler) SetPositionValues(c *gin.Context) {
	position := c.Param("position")
	if !isValidPosition(position) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的位置标识",
		})
		return
	}

	var req struct {
		Values []string `json:"values" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	h.service.SetPositionValues(position, req.Values)
	c.JSON(http.StatusOK, gin.H{
		"message": "设置成功",
	})
}

// DeletePositionValue 删除位置值
func (h *PositionHandler) DeletePositionValue(c *gin.Context) {
	position := c.Param("position")
	value := c.Query("value")

	if !isValidPosition(position) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的位置标识",
		})
		return
	}

	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "值不能为空",
		})
		return
	}

	h.service.DeletePositionValue(position, value)
	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

func isValidPosition(position string) bool {
	validPositions := []string{"a", "b", "c", "d"}
	for _, p := range validPositions {
		if p == position {
			return true
		}
	}
	return false
}

