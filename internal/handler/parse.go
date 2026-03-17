package handler

import (
	"net/http"

	"panflow/internal/service"

	"github.com/gin-gonic/gin"
)

type ParseHandler struct {
	parseSvc  *service.ParseService
	configSvc *service.ConfigService
	tokenSvc  *service.TokenService
}

func NewParseHandler(parseSvc *service.ParseService, configSvc *service.ConfigService, tokenSvc *service.TokenService) *ParseHandler {
	return &ParseHandler{parseSvc: parseSvc, configSvc: configSvc, tokenSvc: tokenSvc}
}

// GET /user/parse/config
func (h *ParseHandler) GetConfig(c *gin.Context) {
	ctx := c.Request.Context()
	Success(c, gin.H{
		"guest_daily_limit": h.configSvc.GetInt(ctx, "guest_daily_limit", 5),
		"svip_daily_limit":  h.configSvc.GetInt(ctx, "svip_daily_limit", 100),
		"vip_count_based":   h.configSvc.GetBool(ctx, "vip_count_based", true),
		"admin_unlimited":   h.configSvc.GetBool(ctx, "admin_unlimited", true),
	})
}

// GET /user/parse/limit
func (h *ParseHandler) GetLimit(c *gin.Context) {
	ctx := c.Request.Context()
	Success(c, gin.H{
		"max_once":            h.configSvc.GetInt(ctx, "max_once", 5),
		"min_single_filesize": h.configSvc.GetInt(ctx, "min_single_filesize", 0),
		"max_single_filesize": h.configSvc.GetInt(ctx, "max_single_filesize", 0),
		"max_all_filesize":    h.configSvc.GetInt(ctx, "max_all_filesize", 0),
	})
}

// POST /user/parse/get_file_list
func (h *ParseHandler) GetFileList(c *gin.Context) {
	var req struct {
		Surl string `json:"surl" binding:"required"`
		Pwd  string `json:"pwd"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}
	// Placeholder: full implementation in parse service
	Success(c, gin.H{"surl": req.Surl})
}

// POST /user/parse/get_vcode
func (h *ParseHandler) GetVcode(c *gin.Context) {
	// Placeholder: captcha handling
	Success(c, gin.H{"vcode": ""})
}

// POST /user/parse/get_download_links
func (h *ParseHandler) GetDownloadLinks(c *gin.Context) {
	var req struct {
		Surl    string  `json:"surl" binding:"required"`
		Pwd     string  `json:"pwd"`
		FsIDs   []int64 `json:"fs_id" binding:"required,min=1"`
		TokenStr string `json:"token"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		FailBadRequest(c, 40000, err.Error())
		return
	}

	clientIP, _ := c.Get("client_ip")
	fingerprint, _ := c.Get("fingerprint")
	ip, _ := clientIP.(string)
	fp, _ := fingerprint.(string)

	results, err := h.parseSvc.Parse(c.Request.Context(), &service.ParseRequest{
		Surl:        req.Surl,
		Pwd:         req.Pwd,
		FsIDs:       req.FsIDs,
		ClientIP:    ip,
		Fingerprint: fp,
		UA:          c.GetHeader("User-Agent"),
		TokenStr:    req.TokenStr,
	})
	if err != nil {
		Fail(c, http.StatusBadRequest, 40010, err.Error())
		return
	}

	Success(c, results)
}
