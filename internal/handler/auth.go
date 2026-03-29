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
	jwtSvc        *service.JWTService
	userRepo      *repository.UserRepository
	adminPassword string
	refreshTTL    time.Duration
}

func NewAuthHandler(
	jwtSvc *service.JWTService,
	userRepo *repository.UserRepository,
	adminPassword string,
	refreshDays int,
) *AuthHandler {
	return &AuthHandler{
		jwtSvc:        jwtSvc,
		userRepo:      userRepo,
		adminPassword: adminPassword,
		refreshTTL:    time.Duration(refreshDays) * 24 * time.Hour,
	}
}

// POST /admin/login
func (h *AuthHandler) AdminLogin(c *gin.Context) {
	var req struct {
		AdminPassword string `json:"admin_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, "admin_password required")
		return
	}
	if req.AdminPassword != h.adminPassword {
		Fail(c, http.StatusUnauthorized, 20001, "admin password error")
		return
	}

	tokenStr, exp, err := h.jwtSvc.Issue()
	if err != nil {
		FailInternal(c, "failed to issue token")
		return
	}
	Success(c, gin.H{
		"token":      tokenStr,
		"expires_at": exp.Format("2006-01-02 15:04:05"),
	})
}

// POST /user/login — 账号密码登录
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, "username and password required")
		return
	}
	h.loginWithPassword(c, req.Username, req.Password)
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

// ── internal helpers ──────────────────────────────────────────────────────────

func (h *AuthHandler) loginWithPassword(c *gin.Context, username, password string) {
	if password == "" {
		Fail(c, http.StatusBadRequest, 40000, "password required")
		return
	}
	user, err := h.userRepo.GetByUsername(username)
	if err != nil || user.Password == "" {
		Fail(c, http.StatusUnauthorized, 20005, "username or password error")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		Fail(c, http.StatusUnauthorized, 20005, "username or password error")
		return
	}
	h.issueTokenPair(c, 0, user.UserType, &user.ID)
}

// issueTokenPair 生成 access token + refresh token 并写 Redis
func (h *AuthHandler) issueTokenPair(c *gin.Context, tokenID uint, userType string, userID *uint) {
	ctx := c.Request.Context()

	accessStr, exp, err := h.jwtSvc.IssueUser(tokenID, userType, userID)
	if err != nil {
		FailInternal(c, "failed to issue access token")
		return
	}

	refreshStr, err := service.NewRefreshToken()
	if err != nil {
		FailInternal(c, "failed to generate refresh token")
		return
	}

	if err := service.SetRefreshToken(ctx, refreshStr, service.RefreshPayload{
		TokenID:  tokenID,
		UserType: userType,
		UserID:   userID,
	}, h.refreshTTL); err != nil {
		FailInternal(c, "failed to store refresh token")
		return
	}

	Success(c, gin.H{
		"access_token":  accessStr,
		"refresh_token": refreshStr,
		"user_type":     userType,
		"expires_at":    exp.Format("2006-01-02 15:04:05"),
	})
}
