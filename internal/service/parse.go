package service

import (
	"context"
	"errors"
	"fmt"

	"panflow/internal/model"
	"panflow/internal/repository"
)

var (
	ErrFileSizeTooSmall  = errors.New("file size too small")
	ErrFileSizeTooLarge  = errors.New("file size too large")
	ErrTotalSizeTooLarge = errors.New("total size too large")
	ErrNoDownloadURL     = errors.New("no download url returned")
)

// ParseRequest 是解析操作的输入。
// TokenID 和 UserType 来自 JWT 上下文，不在请求体中传递明文 token。
type ParseRequest struct {
	Surl        string
	Pwd         string
	FsIDs       []int64
	ClientIP    string
	Fingerprint string
	UA          string
	TokenID     uint   // from JWT claim
	UserType    string // from JWT claim: guest | vip | svip | admin
	UserID      *uint  // from JWT claim (optional, for svip)
}

// ParseResult 是解析操作的输出
type ParseResult struct {
	FsID int64    `json:"fs_id"`
	URLs []string `json:"urls"`
	Size int64    `json:"size"`
}

// parseCreds 封装各账号类型的认证凭据
type parseCreds struct {
	transferCookie string // 用于 GetShareInfo / GetFileList / TransferFiles
	locateCookie   string // 用于 LocateDownload（download_ticket 时与 transferCookie 不同）
	accessToken    string // open_platform 专用
	ua             string
}

// ParseService 编排完整解析流程
type ParseService struct {
	tokenSvc   *TokenService
	userSvc    *UserService
	accountSvc *AccountService
	recordSvc  *RecordService
	configSvc  *ConfigService
	fileRepo   *repository.FileListRepository
	client     *bdwpClient
	userAgent  string
}

func NewParseService(
	tokenSvc *TokenService,
	userSvc *UserService,
	accountSvc *AccountService,
	recordSvc *RecordService,
	configSvc *ConfigService,
	fileRepo *repository.FileListRepository,
	proxyURL string,
	userAgent string,
) *ParseService {
	return &ParseService{
		tokenSvc:   tokenSvc,
		userSvc:    userSvc,
		accountSvc: accountSvc,
		recordSvc:  recordSvc,
		configSvc:  configSvc,
		fileRepo:   fileRepo,
		client:     newBdwpClient(proxyURL),
		userAgent:  userAgent,
	}
}

// Parse 执行完整的下载链接解析流程
func (s *ParseService) Parse(ctx context.Context, req *ParseRequest) ([]*ParseResult, error) {
	// 1. 从配置加载限制
	minSize := int64(s.configSvc.GetInt(ctx, "min_single_filesize", 0))
	maxSize := int64(s.configSvc.GetInt(ctx, "max_single_filesize", 0))
	maxTotal := int64(s.configSvc.GetInt(ctx, "max_all_filesize", 0))

	// 2. 校验 token 配额（通过 JWT 已认证，直接按 ID 查）
	token, err := s.tokenSvc.ValidateByID(ctx, req.TokenID, req.ClientIP)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	// 3. 若有关联用户则校验用户配额
	var user *model.User
	if req.UserID != nil {
		user, err = s.userSvc.CheckQuota(ctx, *req.UserID)
		if err != nil {
			return nil, fmt.Errorf("user: %w", err)
		}
	}

	// 4. 选取账号
	account, err := s.accountSvc.PickForUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 5. 按账号类型提取认证凭据
	creds := s.extractCreds(account)
	if creds.ua == "" {
		creds.ua = s.userAgent
	}

	// 6. 获取分享链接元数据
	shareInfo, err := s.client.GetShareInfo(req.Surl, req.Pwd, creds.transferCookie, creds.accessToken, creds.ua)
	if err != nil {
		return nil, fmt.Errorf("share info: %w", err)
	}

	// 7. 获取文件列表以解析文件大小
	fileListResp, err := s.client.GetFileList(req.Surl, req.Pwd, creds.transferCookie, creds.accessToken, creds.ua,
		shareInfo.ShareID, shareInfo.UK, shareInfo.BDSToken)
	if err != nil {
		return nil, fmt.Errorf("file list: %w", err)
	}

	fileMap := make(map[int64]*ShareFile)
	for i := range fileListResp.List {
		f := &fileListResp.List[i]
		fileMap[f.FsID] = f
	}

	// 8. 校验文件大小
	var totalSize int64
	for _, fsID := range req.FsIDs {
		f, ok := fileMap[fsID]
		if !ok {
			continue
		}
		if minSize > 0 && f.Size < minSize {
			return nil, ErrFileSizeTooSmall
		}
		if maxSize > 0 && f.Size > maxSize {
			return nil, ErrFileSizeTooLarge
		}
		totalSize += f.Size
	}
	if maxTotal > 0 && totalSize > maxTotal {
		return nil, ErrTotalSizeTooLarge
	}

	// 9. 转存文件到账号「我的资源」
	if err := s.client.TransferFiles(req.Surl, req.Pwd, req.FsIDs,
		shareInfo.ShareID, shareInfo.UK, shareInfo.BDSToken,
		creds.transferCookie, creds.accessToken, creds.ua); err != nil {
		return nil, fmt.Errorf("transfer: %w", err)
	}

	// 10. 获取高速下载链接，并收集转存路径用于后续清理
	var results []*ParseResult
	var transferredPaths []string
	for _, fsID := range req.FsIDs {
		urls, err := s.client.LocateDownload(fsID, creds.locateCookie, creds.accessToken, creds.ua)
		if err != nil || len(urls) == 0 {
			continue
		}

		size := int64(0)
		if f, ok := fileMap[fsID]; ok {
			size = f.Size
			transferredPaths = append(transferredPaths, "/我的资源/"+f.Filename)
		}

		results = append(results, &ParseResult{FsID: fsID, URLs: urls, Size: size})

		if f, ok := fileMap[fsID]; ok {
			_ = s.fileRepo.Upsert(&model.FileList{
				Surl:     req.Surl,
				Pwd:      req.Pwd,
				FsID:     fmt.Sprintf("%d", fsID),
				Size:     f.Size,
				Filename: f.Filename,
			})
		}
	}

	if len(results) == 0 {
		return nil, ErrNoDownloadURL
	}

	// 11. Method A：LocateDownload 完成后立即删除转存文件
	// CDN 链接已独立于源文件，删除不影响下载
	_ = s.client.DeleteFiles(transferredPaths, creds.transferCookie, creds.accessToken, creds.ua)

	// 12. 记录用量
	_ = s.tokenSvc.RecordUsage(ctx, token.ID, totalSize)
	_ = s.accountSvc.RecordUsage(ctx, account.ID, totalSize)
	if user != nil {
		_ = s.userSvc.RecordUsage(ctx, user.ID, user.UserType)
	}

	// 13. 保存解析记录
	urlStrs := make([]string, 0, len(results))
	for _, r := range results {
		urlStrs = append(urlStrs, r.URLs...)
	}
	record := &model.Record{
		IP:          req.ClientIP,
		Fingerprint: req.Fingerprint,
		UA:          req.UA,
		TokenID:     token.ID,
		AccountID:   account.ID,
		URLs:        model.JSONSlice(urlStrs),
	}
	if req.UserID != nil {
		record.UserID = req.UserID
	}
	_ = s.recordSvc.Save(ctx, record)

	return results, nil
}

// extractCreds 按账号类型提取认证凭据
func (s *ParseService) extractCreds(account *model.Account) parseCreds {
	data := account.AccountData
	var creds parseCreds

	if v, ok := data["user_agent"].(string); ok {
		creds.ua = v
	}

	switch account.AccountType {
	case "open_platform":
		creds.accessToken, _ = data["access_token"].(string)

	case "download_ticket":
		creds.transferCookie, _ = data["save_cookie"].(string)
		creds.locateCookie, _ = data["download_cookie"].(string)

	case "enterprise_cookie":
		creds.transferCookie, _ = data["cookie"].(string)
		if dc, ok := data["dlink_cookie"].(string); ok && dc != "" {
			creds.locateCookie = dc
		} else {
			creds.locateCookie = creds.transferCookie
		}

	default: // cookie
		creds.transferCookie, _ = data["cookie"].(string)
		creds.locateCookie = creds.transferCookie
	}

	return creds
}
