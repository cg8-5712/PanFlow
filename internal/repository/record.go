package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type RecordRepository struct {
	db *gorm.DB
}

func NewRecordRepository(db *gorm.DB) *RecordRepository {
	return &RecordRepository{db: db}
}

func (r *RecordRepository) Create(record *model.Record) error {
	return r.db.Create(record).Error
}

func (r *RecordRepository) List(offset, limit int) ([]model.Record, int64, error) {
	var records []model.Record
	var total int64
	if err := r.db.Model(&model.Record{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Preload("Token").Preload("Account").Preload("User").
		Order("id DESC").Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

func (r *RecordRepository) ListByUserID(userID uint, offset, limit int) ([]model.Record, int64, error) {
	var records []model.Record
	var total int64
	if err := r.db.Model(&model.Record{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Where("user_id = ?", userID).
		Order("id DESC").Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

func (r *RecordRepository) ListByTokenID(tokenID uint, offset, limit int) ([]model.Record, int64, error) {
	var records []model.Record
	var total int64
	if err := r.db.Model(&model.Record{}).Where("token_id = ?", tokenID).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Where("token_id = ?", tokenID).
		Order("id DESC").Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}
