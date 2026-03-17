package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type TokenRepository struct {
	db *gorm.DB
}

func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

func (r *TokenRepository) Create(token *model.Token) error {
	return r.db.Create(token).Error
}

func (r *TokenRepository) GetByID(id uint) (*model.Token, error) {
	var token model.Token
	err := r.db.First(&token, id).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) GetByToken(tokenStr string) (*model.Token, error) {
	var token model.Token
	err := r.db.Where("token = ?", tokenStr).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (r *TokenRepository) Update(token *model.Token) error {
	return r.db.Save(token).Error
}

func (r *TokenRepository) Delete(id uint) error {
	return r.db.Delete(&model.Token{}, id).Error
}

func (r *TokenRepository) List(offset, limit int) ([]model.Token, int64, error) {
	var tokens []model.Token
	var total int64
	if err := r.db.Model(&model.Token{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&tokens).Error
	return tokens, total, err
}

func (r *TokenRepository) IncrementUsage(id uint, size int64) error {
	return r.db.Model(&model.Token{}).Where("id = ?", id).Updates(map[string]any{
		"used_count": gorm.Expr("used_count + 1"),
		"used_size":  gorm.Expr("used_size + ?", size),
	}).Error
}
