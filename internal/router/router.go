package router

import (
	"panflow/internal/config"
	"panflow/internal/handler"
	"panflow/internal/middleware"
	"panflow/internal/repository"
	"panflow/internal/service"

	"github.com/gin-gonic/gin"
)

// Setup registers all routes on the given engine
func Setup(
	r *gin.Engine,
	cfg *config.Config,
	accountRepo *repository.AccountRepository,
	tokenRepo *repository.TokenRepository,
	userRepo *repository.UserRepository,
	configRepo *repository.ConfigRepository,
	recordRepo *repository.RecordRepository,
	fileListRepo *repository.FileListRepository,
	blackListRepo *repository.BlackListRepository,
) {
	// Services
	tokenSvc := service.NewTokenService(tokenRepo)
	userSvc := service.NewUserService(userRepo)
	accountSvc := service.NewAccountService(accountRepo, cfg.Panflow.ProxyHTTP)
	recordSvc := service.NewRecordService(recordRepo)
	configSvc := service.NewConfigService(configRepo)
	jwtSvc := service.NewJWTService(cfg.Panflow.JWTSecret, cfg.Panflow.JWTExpireHours)
	parseSvc := service.NewParseService(
		tokenSvc, userSvc, accountSvc, recordSvc, configSvc,
		fileListRepo,
		cfg.Panflow.ProxyHTTP,
		cfg.Panflow.GuestUserAgent,
	)

	// Handlers
	accountH := handler.NewAccountHandler(accountRepo, accountSvc)
	tokenH := handler.NewTokenHandler(tokenRepo)
	userH := handler.NewUserHandler(userRepo)
	configH := handler.NewConfigHandler(configRepo)
	recordH := handler.NewRecordHandler(recordRepo)
	blackListH := handler.NewBlackListHandler(blackListRepo)
	parseH := handler.NewParseHandler(parseSvc, configSvc)
	authH := handler.NewAuthHandler(jwtSvc, userRepo, cfg.Panflow.AdminPassword, cfg.Panflow.JWTRefreshDays)

	r.Use(middleware.Cors())

	api := r.Group("/api/v1")

	// ── Public routes ─────────────────────────────────────────────────────────
	api.POST("/admin/login", authH.AdminLogin)
	api.POST("/user/login", authH.UserLogin)
	api.POST("/user/refresh", authH.RefreshToken)
	api.POST("/user/logout", authH.Logout)

	// ── User routes（公开，仅 IdentifierFilter）────────────────────────────────
	public := api.Group("")
	public.Use(middleware.IdentifierFilter(blackListRepo))
	{
		public.GET("/user/parse/config", parseH.GetConfig)
		public.GET("/user/parse/limit", parseH.GetLimit)
		public.POST("/user/parse/get_vcode", parseH.GetVcode)
	}

	// ── User routes（需登录，IdentifierFilter + UserJWTAuth）─────────────────
	auth := api.Group("")
	auth.Use(middleware.IdentifierFilter(blackListRepo))
	auth.Use(middleware.UserJWTAuth(jwtSvc))
	{
		auth.POST("/user/parse/get_file_list", parseH.GetFileList)
		auth.POST("/user/parse/get_download_links", parseH.GetDownloadLinks)
		auth.GET("/user/history", recordH.UserHistory)
		auth.GET("/user/profile", userH.GetProfile)
		auth.PATCH("/user/profile", userH.UpdateProfile)
		auth.PATCH("/user/password", userH.ChangePassword)
	}

	// ── Admin routes (JWT protected) ──────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.JWTAuth(jwtSvc))
	{
		admin.GET("/account", accountH.List)
		admin.POST("/account", accountH.Create)
		admin.PATCH("/account", accountH.Update)
		admin.DELETE("/account", accountH.Delete)
		admin.POST("/account/update_data", accountH.UpdateData)
		admin.POST("/account/check_ban_status", accountH.CheckBanStatus)
		admin.POST("/account/check_enterprise_cid", accountH.CheckEnterpriseCID)

		admin.GET("/token", tokenH.List)
		admin.POST("/token", tokenH.Create)
		admin.PATCH("/token", tokenH.Update)
		admin.DELETE("/token", tokenH.Delete)

		admin.GET("/user", userH.List)
		admin.POST("/user", userH.Create)
		admin.PATCH("/user", userH.Update)
		admin.DELETE("/user", userH.Delete)
		admin.POST("/user/recharge", userH.Recharge)
		admin.POST("/user/reset_password", userH.ResetPassword)

		admin.GET("/config", configH.List)
		admin.PATCH("/config", configH.Update)
		admin.POST("/config/reload", configH.Reload)

		admin.GET("/black_list", blackListH.List)
		admin.POST("/black_list", blackListH.Create)
		admin.PATCH("/black_list", blackListH.Update)
		admin.DELETE("/black_list", blackListH.Delete)

		admin.GET("/record", recordH.List)
		admin.GET("/record/history", recordH.History)

	}
}
