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
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if err := h.repo.Create(&user); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, user)
}

// PATCH /admin/user
func (h *UserHandler) Update(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	if user.ID == 0 {
		FailBadRequest(c, 40001, "id required")
		return
	}
	if err := h.repo.Update(&user); err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, user)
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
