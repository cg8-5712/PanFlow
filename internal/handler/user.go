package handler

import (
	"net/http"

	"panflow/internal/model"
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// ── User self-service ─────────────────────────────────────────────────────────

// GET /user/profile
func (h *UserHandler) GetProfile(c *gin.Context) {
	uid, ok := c.Get("jwt_user_id")
	if !ok {
		// token-only user (no user record), return token-level info from JWT
		Success(c, gin.H{
			"user_type": c.GetString("jwt_user_type"),
		})
		return
	}

	user, err := h.repo.GetByID(uid.(uint))
	if err != nil {
		FailNotFound(c, 40400, "user not found")
		return
	}

	remaining := int64(user.DailyLimit) - user.DailyUsedCount
	if remaining < 0 {
		remaining = 0
	}

	Success(c, gin.H{
		"username":         user.Username,
		"email":            user.Email,
		"user_type":        user.UserType,
		"vip_balance":      user.VipBalance,
		"daily_used_count": user.DailyUsedCount,
		"daily_limit":      user.DailyLimit,
		"daily_remaining":  remaining,
	})
}

// PATCH /user/profile
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	uid, ok := c.Get("jwt_user_id")
	if !ok {
		FailForbidden(c, 40301, "no user account linked")
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}

	if err := h.repo.UpdateEmail(uid.(uint), req.Email); err != nil {
		FailInternal(c, err.Error())
		return
	}
	SuccessMsg(c, "updated")
}

// ── Admin management ──────────────────────────────────────────────────────────

// GET /admin/user
func (h *UserHandler) List(c *gin.Context) {
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

	users, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": users})
}

// POST /admin/user
func (h *UserHandler) Create(c *gin.Context) {
	var req struct {
		Username   string `json:"username" binding:"required"`
		Email      string `json:"email"`
		UserType   string `json:"user_type"`
		DailyLimit int    `json:"daily_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if req.UserType == "" {
		req.UserType = "guest"
	}
	if req.DailyLimit == 0 {
		req.DailyLimit = 5
	}

	user := model.User{
		Username:   req.Username,
		Email:      req.Email,
		UserType:   req.UserType,
		DailyLimit: req.DailyLimit,
	}
	if err := h.repo.Create(&user); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, user)
}

// PATCH /admin/user
func (h *UserHandler) Update(c *gin.Context) {
	var req struct {
		ID         uint   `json:"id" binding:"required"`
		Email      string `json:"email"`
		UserType   string `json:"user_type"`
		DailyLimit *int   `json:"daily_limit"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}

	fields := map[string]any{}
	if req.Email != "" {
		fields["email"] = req.Email
	}
	if req.UserType != "" {
		fields["user_type"] = req.UserType
	}
	if req.DailyLimit != nil {
		fields["daily_limit"] = *req.DailyLimit
	}
	if len(fields) == 0 {
		FailBadRequest(c, 40001, "no fields to update")
		return
	}

	if err := h.repo.UpdateFields(req.ID, fields); err != nil {
		FailInternal(c, err.Error())
		return
	}
	SuccessMsg(c, "updated")
}

// DELETE /admin/user
func (h *UserHandler) Delete(c *gin.Context) {
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

// POST /admin/user/recharge
func (h *UserHandler) Recharge(c *gin.Context) {
	var req struct {
		ID    uint  `json:"id" binding:"required"`
		Count int64 `json:"count" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.AddVipBalance(req.ID, req.Count); err != nil {
		FailInternal(c, err.Error())
		return
	}
	SuccessMsg(c, "recharged")
}
