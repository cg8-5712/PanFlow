package repository

import (
	"time"

	"panflow/internal/model"

	"gorm.io/gorm"
)

type BlackListRepository struct {
	db *gorm.DB
}

func NewBlackListRepository(db *gorm.DB) *BlackListRepository {
	return &BlackListRepository{db: db}
}

func (r *BlackListRepository) Create(bl *model.BlackList) error {
	return r.db.Create(bl).Error
}

func (r *BlackListRepository) GetByTypeAndIdentifier(typ, identifier string) (*model.BlackList, error) {
	var bl model.BlackList
	err := r.db.Where("type = ? AND identifier = ?", typ, identifier).First(&bl).Error
	if err != nil {
		return nil, err
	}
	return &bl, nil
}

// IsBlocked checks if an identifier is currently blocked (not expired)
func (r *BlackListRepository) IsBlocked(typ, identifier string) (bool, error) {
	var count int64
	err := r.db.Model(&model.BlackList{}).
		Where("type = ? AND identifier = ? AND (expires_at IS NULL OR expires_at > ?)", typ, identifier, time.Now()).
		Count(&count).Error
	return count > 0, err
}

func (r *BlackListRepository) Update(bl *model.BlackList) error {
	return r.db.Save(bl).Error
}

func (r *BlackListRepository) Delete(id uint) error {
	return r.db.Delete(&model.BlackList{}, id).Error
}

func (r *BlackListRepository) List(offset, limit int) ([]model.BlackList, int64, error) {
	var list []model.BlackList
	var total int64
	if err := r.db.Model(&model.BlackList{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&list).Error
	return list, total, err
}

// DeleteExpired removes all expired blacklist entries
func (r *BlackListRepository) DeleteExpired() error {
	return r.db.Where("expires_at IS NOT NULL AND expires_at <= ?", time.Now()).
		Delete(&model.BlackList{}).Error
}
