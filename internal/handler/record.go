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

// GET /admin/record/history  (alias with token filter)
func (h *RecordHandler) History(c *gin.Context) {
	var q struct {
		Page    int  `form:"page"`
		Limit   int  `form:"limit"`
		TokenID uint `form:"token_id"`
		UserID  uint `form:"user_id"`
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

	if q.TokenID > 0 {
		records, total, err := h.repo.ListByTokenID(q.TokenID, offset, q.Limit)
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

	// token_id from context (set by token middleware later)
	tokenID, _ := c.Get("token_id")
	tid, _ := tokenID.(uint)

	records, total, err := h.repo.ListByTokenID(tid, (q.Page-1)*q.Limit, q.Limit)
	if err != nil {
		FailInternal(c, err.Error())
		return
	}
	Success(c, gin.H{"total": total, "list": records})
}
