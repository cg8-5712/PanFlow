package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type ProxyRepository struct {
	db *gorm.DB
}

func NewProxyRepository(db *gorm.DB) *ProxyRepository {
	return &ProxyRepository{db: db}
}

func (r *ProxyRepository) Create(proxy *model.Proxy) error {
	return r.db.Create(proxy).Error
}

func (r *ProxyRepository) GetByID(id uint) (*model.Proxy, error) {
	var proxy model.Proxy
	err := r.db.First(&proxy, id).Error
	if err != nil {
		return nil, err
	}
	return &proxy, nil
}

func (r *ProxyRepository) Update(proxy *model.Proxy) error {
	return r.db.Save(proxy).Error
}

func (r *ProxyRepository) Delete(id uint) error {
	return r.db.Delete(&model.Proxy{}, id).Error
}

func (r *ProxyRepository) List(offset, limit int) ([]model.Proxy, int64, error) {
	var proxies []model.Proxy
	var total int64
	if err := r.db.Model(&model.Proxy{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&proxies).Error
	return proxies, total, err
}

// ListEnabled returns all enabled proxies, optionally filtered by account
func (r *ProxyRepository) ListEnabled(accountID uint) ([]model.Proxy, error) {
	var proxies []model.Proxy
	q := r.db.Where("enable = ?", true)
	if accountID > 0 {
		q = q.Where("account_id = ?", accountID)
	}
	err := q.Find(&proxies).Error
	return proxies, err
}
