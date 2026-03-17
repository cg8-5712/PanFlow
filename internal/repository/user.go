package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create creates a new user
func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

// Delete soft deletes a user
func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

// List retrieves all users with pagination
func (r *UserRepository) List(offset, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

// IncrementDailyUsed increments the daily used count for a user
func (r *UserRepository) IncrementDailyUsed(id uint, count int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("daily_used_count", gorm.Expr("daily_used_count + ?", count)).Error
}

// ResetDailyUsed resets the daily used count for all users
func (r *UserRepository) ResetDailyUsed() error {
	return r.db.Model(&model.User{}).Where("daily_used_count > 0").
		Update("daily_used_count", 0).Error
}

// DeductVipBalance deducts VIP balance for a user
func (r *UserRepository) DeductVipBalance(id uint, count int64) error {
	return r.db.Model(&model.User{}).Where("id = ? AND vip_balance >= ?", id, count).
		UpdateColumn("vip_balance", gorm.Expr("vip_balance - ?", count)).Error
}

// AddVipBalance adds VIP balance for a user (recharge)
func (r *UserRepository) AddVipBalance(id uint, count int64) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).
		UpdateColumn("vip_balance", gorm.Expr("vip_balance + ?", count)).Error
}
