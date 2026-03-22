package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"
	"panflow/internal/service"

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
	repo       *repository.AccountRepository
	accountSvc *service.AccountService
}

func NewAccountHandler(repo *repository.AccountRepository, accountSvc *service.AccountService) *AccountHandler {
	return &AccountHandler{repo: repo, accountSvc: accountSvc}
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

	accounts, total, err := h.repo.ListWithTodayStats((q.Page-1)*q.Limit, q.Limit)
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

// POST /admin/account/update_data
func (h *AccountHandler) UpdateData(c *gin.Context) {
	var req struct {
		ID   uint            `json:"id" binding:"required"`
		Data model.JSONMap   `json:"data" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.UpdateData(req.ID, req.Data); err != nil {
		FailInternal(c, err.Error())
		return
	}
	SuccessMsg(c, "updated")
}

// POST /admin/account/check_ban_status
func (h *AccountHandler) CheckBanStatus(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	acc, err := h.repo.GetByID(req.ID)
	if err != nil {
		Fail(c, http.StatusNotFound, 40400, "account not found")
		return
	}

	var cookieOrToken, accountType string
	var cid int64
	accountType = acc.AccountType

	switch acc.AccountType {
	case "cookie":
		cookieOrToken, _ = acc.AccountData["cookie"].(string)
	case "enterprise_cookie":
		cookieOrToken, _ = acc.AccountData["cookie"].(string)
		if v, ok := acc.AccountData["cid"].(float64); ok {
			cid = int64(v)
		}
	case "open_platform":
		cookieOrToken, _ = acc.AccountData["access_token"].(string)
	case "download_ticket":
		cookieOrToken, _ = acc.AccountData["download_cookie"].(string)
	default:
		FailBadRequest(c, 40001, "unsupported account type for ban check")
		return
	}

	if cookieOrToken == "" {
		FailBadRequest(c, 40001, "account has no cookie or access_token")
		return
	}

	ua := "netdisk;P2SP;3.0.20.138"
	status, err := h.accountSvc.CheckBanStatus(accountType, cookieOrToken, ua, cid)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, status)
}

// POST /admin/account/check_enterprise_cid
func (h *AccountHandler) CheckEnterpriseCID(c *gin.Context) {
	var req struct {
		ID uint `json:"id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	acc, err := h.repo.GetByID(req.ID)
	if err != nil {
		Fail(c, http.StatusNotFound, 40400, "account not found")
		return
	}
	if acc.AccountType != "enterprise_cookie" {
		FailBadRequest(c, 40001, "account is not enterprise_cookie type")
		return
	}

	cookie, _ := acc.AccountData["cookie"].(string)
	if cookie == "" {
		FailBadRequest(c, 40001, "account has no cookie")
		return
	}

	ua := "netdisk;P2SP;3.0.20.138"
	actualCID, err := h.accountSvc.GetEnterpriseCID(cookie, ua)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}

	storedCID, _ := acc.AccountData["cid"].(float64)
	match := int64(storedCID) == actualCID
	Success(c, gin.H{
		"actual_cid": actualCID,
		"stored_cid": int64(storedCID),
		"match":      match,
	})
}
