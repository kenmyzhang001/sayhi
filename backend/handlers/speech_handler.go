package handlers

import (
	"net/http"
	"sayhi/backend/models"
	"sayhi/backend/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SpeechHandler 话术处理器
type SpeechHandler struct {
	service *services.SpeechService
}

// NewSpeechHandler 创建话术处理器
func NewSpeechHandler(service *services.SpeechService) *SpeechHandler {
	return &SpeechHandler{
		service: service,
	}
}

// CreateGroup 创建话术组
func (h *SpeechHandler) CreateGroup(c *gin.Context) {
	var req models.SpeechGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	group, err := h.service.CreateGroup(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetGroup 获取话术组
func (h *SpeechHandler) GetGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	group, err := h.service.GetGroup(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// GetAllGroups 获取所有话术组
func (h *SpeechHandler) GetAllGroups(c *gin.Context) {
	groups := h.service.GetAllGroups()
	c.JSON(http.StatusOK, models.SpeechGroupListResponse{
		Groups: groups,
		Total:  len(groups),
	})
}

// UpdateGroup 更新话术组
func (h *SpeechHandler) UpdateGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	var req models.SpeechGroupUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求参数错误: " + err.Error(),
		})
		return
	}

	group, err := h.service.UpdateGroup(id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, group)
}

// DeleteGroup 删除话术组
func (h *SpeechHandler) DeleteGroup(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "无效的ID",
		})
		return
	}

	if err := h.service.DeleteGroup(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

