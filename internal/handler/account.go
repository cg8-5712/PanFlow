package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	repo *repository.ConfigRepository
}

func NewConfigHandler(repo *repository.ConfigRepository) *ConfigHandler {
	return &ConfigHandler{repo: repo}
}

// GET /admin/config
func (h *ConfigHandler) List(c *gin.Context) {
	configs, err := h.repo.GetAll()
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, configs)
}

// PATCH /admin/config
func (h *ConfigHandler) Update(c *gin.Context) {
	var req struct {
		Key         string `json:"key" binding:"required"`
		Value       string `json:"value" binding:"required"`
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}

	if err := h.repo.Set(req.Key, req.Value, req.Type, req.Description); err != nil {
		FailInternal(c, err.Error())
		return
	}
	SuccessMsg(c, "updated")
}

// POST /admin/config/reload  (placeholder — services will hook into this)
func (h *ConfigHandler) Reload(c *gin.Context) {
	SuccessMsg(c, "reloaded")
}

type AccountHandler struct {
	repo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

// GET /admin/account
func (h *AccountHandler) List(c *gin.Context) {
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

	accounts, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": accounts})
}

// POST /admin/account
func (h *AccountHandler) Create(c *gin.Context) {
	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Create(&account); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, account)
}

// PATCH /admin/account
func (h *AccountHandler) Update(c *gin.Context) {
	var account model.Account
	if err := c.ShouldBindJSON(&account); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if account.ID == 0 {
		FailBadRequest(c, 40001, "id required")
		return
	}
	if err := h.repo.Update(&account); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, account)
}

// DELETE /admin/account
func (h *AccountHandler) Delete(c *gin.Context) {
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
