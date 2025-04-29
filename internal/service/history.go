package service

import (
	"barcode-checker/internal/model"
	"barcode-checker/internal/repository"
)

type HistoryService interface {
	GetUserHistory(userID uint, page, limit int) ([]model.ProductCheck, int64, error)
	DeleteHistoryItem(userID uint, id string) error
}

type historyService struct {
	repo repository.HistoryRepository
}

func NewHistoryService(repo repository.HistoryRepository) HistoryService {
	return &historyService{repo: repo}
}

func (s *historyService) GetUserHistory(userID uint, page, limit int) ([]model.ProductCheck, int64, error) {
	return s.repo.GetByUserID(userID, page, limit)
}

func (s *historyService) DeleteHistoryItem(userID uint, id string) error {
	return s.repo.DeleteByID(userID, id)
}
