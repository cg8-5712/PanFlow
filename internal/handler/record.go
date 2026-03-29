package handler

import (
	"panflow/internal/repository"

	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	repo *repository.RecordRepository
}

func NewRecordHandler(repo *repository.RecordRepository) *RecordHandler {
	return &RecordHandler{repo: repo}
}

// GET /admin/record
func (h *RecordHandler) List(c *gin.Context) {
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

	records, total, err := h.repo.List((q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": records})
}

// GET /admin/record/history
func (h *RecordHandler) History(c *gin.Context) {
	var q struct {
		Page   int  `form:"page"`
		Limit  int  `form:"limit"`
		UserID uint `form:"user_id"`
	}
	_ = c.ShouldBindQuery(&q)
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 100 {
		q.Limit = 20
	}

	offset := (q.Page - 1) * q.Limit

	if q.UserID > 0 {
		records, total, err := h.repo.ListByUserID(q.UserID, offset, q.Limit)
		if err != nil {
			FailInternal(c, err.Error())
			return
		}
		Success(c, gin.H{"total": total, "list": records})
		return
	}

	records, total, err := h.repo.List(offset, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": records})
}

// GET /user/history  (user's own history)
func (h *RecordHandler) UserHistory(c *gin.Context) {
	var q struct {
		Page  int `form:"page"`
		Limit int `form:"limit"`
	}
	_ = c.ShouldBindQuery(&q)
	if q.Page < 1 {
		q.Page = 1
	}
	if q.Limit < 1 || q.Limit > 50 {
		q.Limit = 10
	}

	userID, _ := c.Get("jwt_user_id")
	uid, _ := userID.(uint)

	records, total, err := h.repo.ListByUserID(uid, (q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": records})
}
