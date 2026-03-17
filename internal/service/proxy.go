package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"panflow/internal/model"
	"panflow/internal/repository"
)

type ProxyService struct {
	repo *repository.ProxyRepository
}

func NewProxyService(repo *repository.ProxyRepository) *ProxyService {
	return &ProxyService{repo: repo}
}

// PickForAccount returns a random enabled proxy for the given account, or nil if none
func (s *ProxyService) PickForAccount(ctx context.Context, accountID uint) (*model.Proxy, error) {
	proxies, err := s.repo.ListEnabled(accountID)
	if err != nil {
		return nil, err
	}
	if len(proxies) == 0 {
		// fall back to global proxies
		proxies, err = s.repo.ListEnabled(0)
		if err != nil || len(proxies) == 0 {
			return nil, nil
		}
	}
	return &proxies[rand.Intn(len(proxies))], nil
}

// BuildHTTPClient returns an *http.Client configured with the given proxy
func (s *ProxyService) BuildHTTPClient(proxy *model.Proxy) *http.Client {
	transport := &http.Transport{}
	if proxy != nil && proxy.Enable {
		if u, err := url.Parse(proxy.Proxy); err == nil {
			transport.Proxy = http.ProxyURL(u)
		}
	}
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
}

// WrapURL optionally rewrites a download URL through a proxy server
// proxyServer format: "https://proxy.example.com"
func (s *ProxyService) WrapURL(downloadURL, proxyServer, proxyPassword string) string {
	if proxyServer == "" {
		return downloadURL
	}
	encoded := url.QueryEscape(downloadURL)
	return fmt.Sprintf("%s/d?url=%s&password=%s", proxyServer, encoded, url.QueryEscape(proxyPassword))
}
