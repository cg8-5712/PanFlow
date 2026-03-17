package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baiduPCSBase    = "https://pan.baidu.com"
	baiduAPIBase    = "https://pan.baidu.com/api"
	baiduDiskBase   = "https://pan.baidu.com/rest/2.0/xpan"
	locateDownload  = "https://pan.baidu.com/api/locatedownload"
	defaultTimeout  = 30 * time.Second
)

// bdwpClient is a thin HTTP client for Baidu Pan API calls
type bdwpClient struct {
	http    *http.Client
	proxyFn func(*http.Request) (*url.URL, error)
}

func newBdwpClient(proxyURL string) *bdwpClient {
	transport := &http.Transport{}
	if proxyURL != "" {
		if u, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(u)
		}
	}
	return &bdwpClient{
		http: &http.Client{
			Timeout:   defaultTimeout,
			Transport: transport,
		},
	}
}

func (c *bdwpClient) get(rawURL string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (c *bdwpClient) post(rawURL string, form url.Values, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, rawURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// ─── Share link info ──────────────────────────────────────────────────────────

type ShareInfo struct {
	Errno    int    `json:"errno"`
	ShareID  int64  `json:"shareid"`
	UK       int64  `json:"uk"`
	BDSToken string `json:"bdstoken"`
}

// GetShareInfo fetches share link metadata (shareid, uk, bdstoken)
func (c *bdwpClient) GetShareInfo(surl, pwd, cookie, userAgent string) (*ShareInfo, error) {
	apiURL := fmt.Sprintf("%s/share/wxlist?shorturl=%s&root=1", baiduPCSBase, surl)
	if pwd != "" {
		apiURL += "&pwd=" + url.QueryEscape(pwd)
	}

	body, err := c.get(apiURL, map[string]string{
		"Cookie":     cookie,
		"User-Agent": userAgent,
		"Referer":    baiduPCSBase,
	})
	if err != nil {
		return nil, fmt.Errorf("get share info: %w", err)
	}

	var info ShareInfo
	if err := json.Unmarshal(body, &info); err != nil {
		return nil, fmt.Errorf("parse share info: %w", err)
	}
	if info.Errno != 0 {
		return nil, fmt.Errorf("share info errno %d", info.Errno)
	}
	return &info, nil
}

// ─── File list ────────────────────────────────────────────────────────────────

type ShareFile struct {
	FsID     int64  `json:"fs_id"`
	Filename string `json:"server_filename"`
	Size     int64  `json:"size"`
	IsDir    int    `json:"isdir"`
	Path     string `json:"path"`
}

type FileListResp struct {
	Errno int         `json:"errno"`
	List  []ShareFile `json:"list"`
}

// GetFileList fetches the file list from a share link
func (c *bdwpClient) GetFileList(surl, pwd, cookie, userAgent string, shareID, uk int64, bdstoken string) (*FileListResp, error) {
	params := url.Values{}
	params.Set("shorturl", surl)
	params.Set("shareid", fmt.Sprintf("%d", shareID))
	params.Set("uk", fmt.Sprintf("%d", uk))
	params.Set("bdstoken", bdstoken)
	if pwd != "" {
		params.Set("pwd", pwd)
	}

	apiURL := fmt.Sprintf("%s/share/list?%s", baiduPCSBase, params.Encode())
	body, err := c.get(apiURL, map[string]string{
		"Cookie":     cookie,
		"User-Agent": userAgent,
		"Referer":    baiduPCSBase,
	})
	if err != nil {
		return nil, fmt.Errorf("get file list: %w", err)
	}

	var resp FileListResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse file list: %w", err)
	}
	return &resp, nil
}

// ─── Transfer (save to account) ───────────────────────────────────────────────

type TransferResp struct {
	Errno int `json:"errno"`
}

// TransferFiles saves share files into the account's "我的资源" directory
func (c *bdwpClient) TransferFiles(surl, pwd string, fsIDs []int64, shareID, uk int64, bdstoken, cookie, userAgent string) error {
	fsIDsJSON, _ := json.Marshal(fsIDs)

	form := url.Values{}
	form.Set("shorturl", surl)
	form.Set("shareid", fmt.Sprintf("%d", shareID))
	form.Set("from", fmt.Sprintf("%d", uk))
	form.Set("bdstoken", bdstoken)
	form.Set("fsidlist", string(fsIDsJSON))
	form.Set("path", "/我的资源")
	if pwd != "" {
		form.Set("pwd", pwd)
	}

	apiURL := fmt.Sprintf("%s/share/transfer?ondup=newcopy", baiduAPIBase)
	body, err := c.post(apiURL, form, map[string]string{
		"Cookie":     cookie,
		"User-Agent": userAgent,
		"Referer":    fmt.Sprintf("%s/s/%s", baiduPCSBase, surl),
	})
	if err != nil {
		return fmt.Errorf("transfer files: %w", err)
	}

	var resp TransferResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parse transfer resp: %w", err)
	}
	if resp.Errno != 0 {
		return fmt.Errorf("transfer errno %d", resp.Errno)
	}
	return nil
}

// ─── Locate download ──────────────────────────────────────────────────────────

type LocateResp struct {
	Errno int `json:"errno"`
	URLs  []struct {
		URL string `json:"url"`
	} `json:"urls"`
}

// LocateDownload fetches the high-speed download URL for a file
func (c *bdwpClient) LocateDownload(fsID int64, cookie, userAgent string) ([]string, error) {
	params := url.Values{}
	params.Set("method", "locatedownload")
	params.Set("ver", "4.0")
	params.Set("fs_id", fmt.Sprintf("%d", fsID))
	params.Set("path", "/我的资源")

	apiURL := fmt.Sprintf("%s/file?%s", baiduDiskBase, params.Encode())
	body, err := c.get(apiURL, map[string]string{
		"Cookie":     cookie,
		"User-Agent": userAgent,
	})
	if err != nil {
		return nil, fmt.Errorf("locate download: %w", err)
	}

	var resp LocateResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("parse locate resp: %w", err)
	}
	if resp.Errno != 0 {
		return nil, fmt.Errorf("locate errno %d", resp.Errno)
	}

	urls := make([]string, 0, len(resp.URLs))
	for _, u := range resp.URLs {
		if u.URL != "" {
			urls = append(urls, u.URL)
		}
	}
	return urls, nil
}

// ─── Delete file ──────────────────────────────────────────────────────────────

type DeleteResp struct {
	Errno int `json:"errno"`
}

// DeleteFile removes a file from the account's disk
func (c *bdwpClient) DeleteFile(path, bdstoken, cookie, userAgent string) error {
	pathsJSON, _ := json.Marshal([]string{path})

	form := url.Values{}
	form.Set("filelist", string(pathsJSON))
	form.Set("bdstoken", bdstoken)

	apiURL := fmt.Sprintf("%s/file/manager?opera=delete", baiduAPIBase)
	body, err := c.post(apiURL, form, map[string]string{
		"Cookie":     cookie,
		"User-Agent": userAgent,
	})
	if err != nil {
		return fmt.Errorf("delete file: %w", err)
	}

	var resp DeleteResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return fmt.Errorf("parse delete resp: %w", err)
	}
	if resp.Errno != 0 {
		return fmt.Errorf("delete errno %d", resp.Errno)
	}
	return nil
}
