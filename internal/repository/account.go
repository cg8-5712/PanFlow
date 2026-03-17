package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *model.Account) error {
	return r.db.Create(account).Error
}

func (r *AccountRepository) GetByID(id uint) (*model.Account, error) {
	var account model.Account
	err := r.db.First(&account, id).Error
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) Update(account *model.Account) error {
	return r.db.Save(account).Error
}

func (r *AccountRepository) Delete(id uint) error {
	return r.db.Delete(&model.Account{}, id).Error
}

func (r *AccountRepository) List(offset, limit int) ([]model.Account, int64, error) {
	var accounts []model.Account
	var total int64
	if err := r.db.Model(&model.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&accounts).Error
	return accounts, total, err
}

// ListEnabled returns all enabled accounts, optionally filtered by type
func (r *AccountRepository) ListEnabled(accountType string) ([]model.Account, error) {
	var accounts []model.Account
	q := r.db.Where("`switch` = ?", true)
	if accountType != "" {
		q = q.Where("account_type = ?", accountType)
	}
	err := q.Find(&accounts).Error
	return accounts, err
}

// ListByProviderUser returns accounts belonging to a specific user
func (r *AccountRepository) ListByProviderUser(userID uint) ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Where("provider_user_id = ?", userID).Find(&accounts).Error
	return accounts, err
}

func (r *AccountRepository) IncrementUsage(id uint, size int64) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Updates(map[string]any{
		"used_count":  gorm.Expr("used_count + 1"),
		"used_size":   gorm.Expr("used_size + ?", size),
		"last_use_at": gorm.Expr("NOW()"),
	}).Error
}
