package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type BlackListHandler struct {
	repo *repository.BlackListRepository
}

func NewBlackListHandler(repo *repository.BlackListRepository) *BlackListHandler {
	return &BlackListHandler{repo: repo}
}

// GET /admin/black_list
func (h *BlackListHandler) List(c *gin.Context) {
	var q struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}
	_ = c.ShouldBindQuery(&q)
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}

	list, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": list})
}

// POST /admin/black_list
func (h *BlackListHandler) Create(c *gin.Context) {
	var bl model.BlackList
	if err := c.ShouldBindJSON(&bl); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Create(&bl); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, bl)
}

// PATCH /admin/black_list
func (h *BlackListHandler) Update(c *gin.Context) {
	var bl model.BlackList
	if err := c.ShouldBindJSON(&bl); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if bl.ID == 0 {
		FailBadRequest(c, 40001, "id required")
		return
	}
	if err := h.repo.Update(&bl); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, bl)
}

// DELETE /admin/black_list
func (h *BlackListHandler) Delete(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Delete(req.ID); err != nil {
		FailInternal(c, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
