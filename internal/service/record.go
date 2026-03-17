package service

import (
	"context"

	"panflow/internal/model"
	"panflow/internal/repository"
)

type RecordService struct {
	repo *repository.RecordRepository
}

func NewRecordService(repo *repository.RecordRepository) *RecordService {
	return &RecordService{repo: repo}
}

// Save persists a new parse record
func (s *RecordService) Save(ctx context.Context, record *model.Record) error {
	return s.repo.Create(record)
}
