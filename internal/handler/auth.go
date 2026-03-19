package handler

import (
	"net/http"

	"panflow/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	jwtSvc        *service.JWTService
	tokenSvc      *service.TokenService
	adminPassword string
}

func NewAuthHandler(jwtSvc *service.JWTService, tokenSvc *service.TokenService, adminPassword string) *AuthHandler {
	return &AuthHandler{jwtSvc: jwtSvc, tokenSvc: tokenSvc, adminPassword: adminPassword}
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

// POST /user/login
// Accepts a token string (API key), validates it, and issues a user JWT.
// Guests can login with token="guest".
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req struct {
		Token string `json:"token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, "token required")
		return
	}

	tok, err := h.tokenSvc.GetByToken(c.Request.Context(), req.Token)
	if err != nil {
		Fail(c, http.StatusUnauthorized, 20003, "token not found or invalid")
		return
	}

	if !tok.Switch {
		Fail(c, http.StatusForbidden, 20004, "token is disabled")
		return
	}

	// provider_user_id links this token to a user account (used for svip)
	jwtStr, exp, err := h.jwtSvc.IssueUser(tok.ID, tok.UserType, tok.ProviderUserID)
	if err != nil {
		FailInternal(c, "failed to issue token")
		return
	}

	Success(c, gin.H{
		"token":      jwtStr,
		"user_type":  tok.UserType,
		"expires_at": exp.Format("2006-01-02 15:04:05"),
	})
}
