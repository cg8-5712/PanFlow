package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	repo *repository.ProxyRepository
}

func NewProxyHandler(repo *repository.ProxyRepository) *ProxyHandler {
	return &ProxyHandler{repo: repo}
}

// GET /admin/proxy
func (h *ProxyHandler) List(c *gin.Context) {
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

	proxies, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": proxies})
}

// POST /admin/proxy
func (h *ProxyHandler) Create(c *gin.Context) {
	var proxy model.Proxy
	if err := c.ShouldBindJSON(&proxy); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Create(&proxy); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, proxy)
}

// PATCH /admin/proxy
func (h *ProxyHandler) Update(c *gin.Context) {
	var proxy model.Proxy
	if err := c.ShouldBindJSON(&proxy); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if proxy.ID == 0 {
		FailBadRequest(c, 40001, "id required")
		return
	}
	if err := h.repo.Update(&proxy); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, proxy)
}

// DELETE /admin/proxy
func (h *ProxyHandler) Delete(c *gin.Context) {
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
