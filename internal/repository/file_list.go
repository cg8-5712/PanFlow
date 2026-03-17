package repository

import (
	"panflow/internal/model"

	"gorm.io/gorm"
)

type FileListRepository struct {
	db *gorm.DB
}

func NewFileListRepository(db *gorm.DB) *FileListRepository {
	return &FileListRepository{db: db}
}

func (r *FileListRepository) GetByFsID(fsID string) (*model.FileList, error) {
	var f model.FileList
	err := r.db.Where("fs_id = ?", fsID).First(&f).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

func (r *FileListRepository) Upsert(f *model.FileList) error {
	var existing model.FileList
	err := r.db.Where("fs_id = ?", f.FsID).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(f).Error
	}
	if err != nil {
		return err
	}
	f.ID = existing.ID
	return r.db.Save(f).Error
}

func (r *FileListRepository) List(offset, limit int) ([]model.FileList, int64, error) {
	var files []model.FileList
	var total int64
	if err := r.db.Model(&model.FileList{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := r.db.Offset(offset).Limit(limit).Find(&files).Error
	return files, total, err
}
