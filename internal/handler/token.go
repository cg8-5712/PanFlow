package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type TokenHandler struct {
	repo *repository.TokenRepository
}

func NewTokenHandler(repo *repository.TokenRepository) *TokenHandler {
	return &TokenHandler{repo: repo}
}

// GET /admin/token
func (h *TokenHandler) List(c *gin.Context) {
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

	tokens, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": tokens})
}

// POST /admin/token
func (h *TokenHandler) Create(c *gin.Context) {
	var token model.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Create(&token); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, token)
}

// PATCH /admin/token
func (h *TokenHandler) Update(c *gin.Context) {
	var token model.Token
	if err := c.ShouldBindJSON(&token); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if token.ID == 0 {
		FailBadRequest(c, 40001, "id required")
		return
	}
	if err := h.repo.Update(&token); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, token)
}

// DELETE /admin/token
func (h *TokenHandler) Delete(c *gin.Context) {
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
