package router

import (
	"panflow/internal/config"
	"panflow/internal/handler"
	"panflow/internal/middleware"
	"panflow/internal/repository"

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
	blackListRepo *repository.BlackListRepository,
	proxyRepo *repository.ProxyRepository,
) {
	// Handlers
	accountH := handler.NewAccountHandler(accountRepo)
	tokenH := handler.NewTokenHandler(tokenRepo)
	userH := handler.NewUserHandler(userRepo)
	configH := handler.NewConfigHandler(configRepo)
	recordH := handler.NewRecordHandler(recordRepo)
	blackListH := handler.NewBlackListHandler(blackListRepo)
	proxyH := handler.NewProxyHandler(proxyRepo)

	r.Use(middleware.Cors())

	api := r.Group("/api/v1")

	// ── User routes ──────────────────────────────────────────────────────────
	user := api.Group("")
	user.Use(middleware.IdentifierFilter(blackListRepo, cfg.Hklist.Debug))
	{
		user.GET("/user/token", func(c *gin.Context) {
			// placeholder: token lookup by query param
			handler.Success(c, nil)
		})
		user.GET("/user/history", recordH.UserHistory)
	}

	// ── Admin routes ─────────────────────────────────────────────────────────
	admin := api.Group("/admin")
	admin.Use(middleware.PassFilterAdmin(&cfg.Hklist))
	{
		// password check
		admin.POST("/check_password", func(c *gin.Context) {
			handler.SuccessMsg(c, "ok")
		})

		// accounts
		admin.GET("/account", accountH.List)
		admin.POST("/account", accountH.Create)
		admin.PATCH("/account", accountH.Update)
		admin.DELETE("/account", accountH.Delete)

		// tokens
		admin.GET("/token", tokenH.List)
		admin.POST("/token", tokenH.Create)
		admin.PATCH("/token", tokenH.Update)
		admin.DELETE("/token", tokenH.Delete)

		// users
		admin.GET("/user", userH.List)
		admin.POST("/user", userH.Create)
		admin.PATCH("/user", userH.Update)
		admin.DELETE("/user", userH.Delete)
		admin.POST("/user/recharge", userH.Recharge)

		// configs
		admin.GET("/config", configH.List)
		admin.PATCH("/config", configH.Update)
		admin.POST("/config/reload", configH.Reload)

		// black lists
		admin.GET("/black_list", blackListH.List)
		admin.POST("/black_list", blackListH.Create)
		admin.PATCH("/black_list", blackListH.Update)
		admin.DELETE("/black_list", blackListH.Delete)

		// records
		admin.GET("/record", recordH.List)
		admin.GET("/record/history", recordH.History)

		// proxies
		admin.GET("/proxy", proxyH.List)
		admin.POST("/proxy", proxyH.Create)
		admin.PATCH("/proxy", proxyH.Update)
		admin.DELETE("/proxy", proxyH.Delete)
	}
}
