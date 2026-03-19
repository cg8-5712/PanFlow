package handler

import (
	"net/http"

	"panflow/internal/repository"
	"panflow/internal/service"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	jwtSvc        *service.JWTService
	tokenSvc      *service.TokenService
	userRepo      *repository.UserRepository
	tokenRepo     *repository.TokenRepository
	adminPassword string
}

func NewAuthHandler(
	jwtSvc *service.JWTService,
	tokenSvc *service.TokenService,
	userRepo *repository.UserRepository,
	tokenRepo *repository.TokenRepository,
	adminPassword string,
) *AuthHandler {
	return &AuthHandler{
		jwtSvc:        jwtSvc,
		tokenSvc:      tokenSvc,
		userRepo:      userRepo,
		tokenRepo:     tokenRepo,
		adminPassword: adminPassword,
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

// POST /user/login
// 支持两种方式：
//   - token 登录：{"token": "xxx"}
//   - 账号密码登录：{"username": "xxx", "password": "xxx"}
func (h *AuthHandler) UserLogin(c *gin.Context) {
	var req struct {
		Token    string `json:"token"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		Fail(c, http.StatusBadRequest, 40000, "invalid request")
		return
	}

	// 账号密码登录
	if req.Username != "" {
		h.loginWithPassword(c, req.Username, req.Password)
		return
	}

	// token 登录
	if req.Token == "" {
		Fail(c, http.StatusBadRequest, 40000, "token or username+password required")
		return
	}
	h.loginWithToken(c, req.Token)
}

func (h *AuthHandler) loginWithToken(c *gin.Context, tokenStr string) {
	tok, err := h.tokenSvc.GetByToken(c.Request.Context(), tokenStr)
	if err != nil {
		Fail(c, http.StatusUnauthorized, 20003, "token not found or invalid")
		return
	}
	if !tok.Switch {
		Fail(c, http.StatusForbidden, 20004, "token is disabled")
		return
	}

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

	// 找到该用户绑定的 token
	tok, err := h.tokenRepo.GetByProviderUserID(user.ID)
	if err != nil {
		Fail(c, http.StatusUnauthorized, 20006, "no active token linked to this account")
		return
	}

	uid := user.ID
	jwtStr, exp, err := h.jwtSvc.IssueUser(tok.ID, user.UserType, &uid)
	if err != nil {
		FailInternal(c, "failed to issue token")
		return
	}

	Success(c, gin.H{
		"token":      jwtStr,
		"user_type":  user.UserType,
		"expires_at": exp.Format("2006-01-02 15:04:05"),
	})
}
