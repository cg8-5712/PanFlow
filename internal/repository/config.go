package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) *ConfigRepository {
	return &ConfigRepository{db: db}
}

// GetByKey retrieves a config by key
func (r *ConfigRepository) GetByKey(key string) (*model.Config, error) {
	var cfg model.Config
	err := r.db.Where("`key` = ?", key).First(&cfg).Error
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

// GetAll retrieves all configs
func (r *ConfigRepository) GetAll() ([]model.Config, error) {
	var configs []model.Config
	err := r.db.Find(&configs).Error
	return configs, err
}

// Set creates or updates a config entry
func (r *ConfigRepository) Set(key, value, typ, description string) error {
	var cfg model.Config
	err := r.db.Where("`key` = ?", key).First(&cfg).Error
	if err == gorm.ErrRecordNotFound {
		// Create new
		cfg = model.Config{
			Key:         key,
			Value:       value,
			Type:        typ,
			Description: description,
		}
		return r.db.Create(&cfg).Error
	}
	if err != nil {
		return err
	}
	// Update existing
	cfg.Value = value
	if typ != "" {
		cfg.Type = typ
	}
	if description != "" {
		cfg.Description = description
	}
	return r.db.Save(&cfg).Error
}

// Delete deletes a config by key
func (r *ConfigRepository) Delete(key string) error {
	return r.db.Where("`key` = ?", key).Delete(&model.Config{}).Error
}

// BatchGet retrieves multiple configs by keys
func (r *ConfigRepository) BatchGet(keys []string) (map[string]*model.Config, error) {
	var configs []model.Config
	err := r.db.Where("`key` IN ?", keys).Find(&configs).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*model.Config)
	for i := range configs {
		result[configs[i].Key] = &configs[i]
	}
	return result, nil
}
