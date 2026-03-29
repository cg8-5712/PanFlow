package repository

import (
	"time"

	"panflow/internal/model"

	"gorm.io/gorm"
)

// AccountWithStats embeds Account and adds today's usage aggregates
type AccountWithStats struct {
	model.Account
	TodayCount int64 `json:"today_count"`
	TodaySize  int64 `json:"today_size"`
}

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

// UpdateData overwrites the account_data field for a given account
func (r *AccountRepository) UpdateData(id uint, data model.JSONMap) error {
	return r.db.Model(&model.Account{}).Where("id = ?", id).Update("account_data", data).Error
}

// ListWithTodayStats returns paginated accounts with today's count and size aggregated from records
func (r *AccountRepository) ListWithTodayStats(offset, limit int) ([]AccountWithStats, int64, error) {
	var total int64
	if err := r.db.Model(&model.Account{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	today := time.Now().Format("2006-01-02")

	var rows []AccountWithStats
	err := r.db.Model(&model.Account{}).
		Select(`accounts.*,
			COALESCE(SUM(CASE WHEN DATE(r.created_at) = ? THEN 1 ELSE 0 END), 0) AS today_count,
			COALESCE(SUM(CASE WHEN DATE(r.created_at) = ? THEN fl.size ELSE 0 END), 0) AS today_size`,
			today, today).
		Joins("LEFT JOIN records r ON r.account_id = accounts.id").
		Joins("LEFT JOIN file_lists fl ON fl.id = r.fs_id").
		Group("accounts.id").
		Offset(offset).Limit(limit).
		Scan(&rows).Error
	return rows, total, err
}
