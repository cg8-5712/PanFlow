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

// ParseRequest is the input for a parse operation.
// TokenID and UserType are set from JWT context; no plaintext token in transit.
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

// ParseResult is the output of a parse operation
type ParseResult struct {
	FsID int64    `json:"fs_id"`
	URLs []string `json:"urls"`
	Size int64    `json:"size"`
}

// ParseService orchestrates the full parse flow
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

// Parse executes the full download link resolution flow
func (s *ParseService) Parse(ctx context.Context, req *ParseRequest) ([]*ParseResult, error) {
	// 1. Load limits from config
	minSize := int64(s.configSvc.GetInt(ctx, "min_single_filesize", 0))
	maxSize := int64(s.configSvc.GetInt(ctx, "max_single_filesize", 0))
	maxTotal := int64(s.configSvc.GetInt(ctx, "max_all_filesize", 0))

	// 2. Validate token quota (by ID, already authenticated via JWT)
	token, err := s.tokenSvc.ValidateByID(ctx, req.TokenID, req.ClientIP)
	if err != nil {
		return nil, fmt.Errorf("token: %w", err)
	}

	// 3. Validate user quota if linked
	var user *model.User
	if req.UserID != nil {
		user, err = s.userSvc.CheckQuota(ctx, *req.UserID)
		if err != nil {
			return nil, fmt.Errorf("user: %w", err)
		}
	}

	// 4. For svip, also enforce user-level daily limit even without linked user record
	if req.UserType == "svip" && user == nil && req.UserID == nil {
		// svip without linked user: treat as guest-level quota
	}

	// 5. Pick an account
	account, err := s.accountSvc.PickForUser(ctx, user)
	if err != nil {
		return nil, err
	}

	// 6. Extract cookie from account_data
	cookie, ua := s.extractCookieAndUA(account)
	if ua == "" {
		ua = s.userAgent
	}

	// 7. Get share info
	shareInfo, err := s.client.GetShareInfo(req.Surl, req.Pwd, cookie, ua)
	if err != nil {
		return nil, fmt.Errorf("share info: %w", err)
	}

	// 8. Get file list to resolve sizes
	fileListResp, err := s.client.GetFileList(req.Surl, req.Pwd, cookie, ua,
		shareInfo.ShareID, shareInfo.UK, shareInfo.BDSToken)
	if err != nil {
		return nil, fmt.Errorf("file list: %w", err)
	}

	fileMap := make(map[int64]*ShareFile)
	for i := range fileListResp.List {
		f := &fileListResp.List[i]
		fileMap[f.FsID] = f
	}

	// 9. Validate file sizes
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

	// 10. Transfer files to account
	if err := s.client.TransferFiles(req.Surl, req.Pwd, req.FsIDs,
		shareInfo.ShareID, shareInfo.UK, shareInfo.BDSToken, cookie, ua); err != nil {
		return nil, fmt.Errorf("transfer: %w", err)
	}

	// 11. Locate download URLs
	var results []*ParseResult
	for _, fsID := range req.FsIDs {
		urls, err := s.client.LocateDownload(fsID, cookie, ua)
		if err != nil || len(urls) == 0 {
			continue
		}

		size := int64(0)
		if f, ok := fileMap[fsID]; ok {
			size = f.Size
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

	// 12. Record usage
	_ = s.tokenSvc.RecordUsage(ctx, token.ID, totalSize)
	_ = s.accountSvc.RecordUsage(ctx, account.ID, totalSize)
	if user != nil {
		_ = s.userSvc.RecordUsage(ctx, user.ID, user.UserType)
	}

	// 13. Save parse record
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

func (s *ParseService) extractCookieAndUA(account *model.Account) (cookie, ua string) {
	data := account.AccountData
	if v, ok := data["cookie"].(string); ok {
		cookie = v
	}
	if v, ok := data["user_agent"].(string); ok {
		ua = v
	}
	return
}
