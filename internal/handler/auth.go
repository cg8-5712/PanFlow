package handler

import (
	"net/http"
	"time"

	"panflow/internal/repository"
	"panflow/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	jwtSvc     *service.JWTService
	userRepo   *repository.UserRepository
	refreshTTL time.Duration
}

func NewAuthHandler(
	jwtSvc *service.JWTService,
	userRepo *repository.UserRepository,
	refreshDays int,
) *AuthHandler {
	return &AuthHandler{
		jwtSvc:     jwtSvc,
		userRepo:   userRepo,
		refreshTTL: time.Duration(refreshDays) * 24 * time.Hour,
	}
}

// POST /user/login — 统一登录，admin/普通用户都走这里
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, "username and password required")
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil || user.Password == "" {
		Fail(c, http.StatusUnauthorized, 20005, "username or password error")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		Fail(c, http.StatusUnauthorized, 20005, "username or password error")
		return
	}
	h.issueTokenPair(c, user.UserType, &user.ID)
}

// POST /user/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, err.Error())
		return
	}

	ctx := c.Request.Context()
	payload, err := service.GetRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		Fail(c, http.StatusUnauthorized, 20007, "refresh token invalid or expired")
		return
	}

	accessStr, exp, err := h.jwtSvc.IssueUser(payload.TokenID, payload.UserType, payload.UserID)
	if err != nil {
		FailInternal(c, "failed to issue token")
		return
	}

	Success(c, gin.H{
		"access_token": accessStr,
		"expires_at":   exp.Format("2006-01-02 15:04:05"),
	})
}

// POST /user/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, err.Error())
		return
	}
	service.DeleteRefreshToken(c.Request.Context(), req.RefreshToken)
	SuccessMsg(c, "logged out")
}

// issueTokenPair 生成 access token + refresh token
func (h *AuthHandler) issueTokenPair(c *gin.Context, userType string, userID *uint) {
	ctx := c.Request.Context()

	accessStr, exp, err := h.jwtSvc.IssueUser(0, userType, userID)
	if err != nil {
		FailInternal(c, "failed to issue access token")
		return
	}

	refreshStr, err := service.NewRefreshToken()
	if err == nil {
		if setErr := service.SetRefreshToken(ctx, refreshStr, service.RefreshPayload{
			UserType: userType,
			UserID:   userID,
		}, h.refreshTTL); setErr != nil {
			refreshStr = "" // Redis 不可用，降级不返回 refresh token
		}
	} else {
		refreshStr = ""
	}

	resp := gin.H{
		"access_token": accessStr,
		"user_type":    userType,
		"expires_at":   exp.Format("2006-01-02 15:04:05"),
	}
	if refreshStr != "" {
		resp["refresh_token"] = refreshStr
	}
	Success(c, resp)
}
