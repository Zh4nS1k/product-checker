package service

import (
	"barcode-checker/internal/model"
	"barcode-checker/internal/repository"
	"barcode-checker/internal/utils"
	"errors"
)

type HistoryService interface {
	GetUserHistory(userID uint, page, limit int) ([]model.ProductCheck, int64, error)
	DeleteHistoryItem(userID uint, id string) error
	UpdateBarcode(userID uint, id string, newBarcode string) error
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

func (s *historyService) UpdateBarcode(userID uint, id string, newBarcode string) error {
	if !utils.IsBarcodeValid(newBarcode) {
		return errors.New("invalid barcode format: must contain only digits and be 8-20 characters long")
	}
	return s.repo.UpdateBarcode(userID, id, newBarcode)
}
