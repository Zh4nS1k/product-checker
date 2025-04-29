package service

import (
	"barcode-checker/internal/model"
	"barcode-checker/internal/repository"
	"barcode-checker/internal/utils"
	"time"
)

type ProductService interface {
	CheckProduct(barcode string, userID uint, clientIP string) (*model.ProductCheckResult, error)
}

type productService struct {
	barcodeChecker BarcodeChecker
	historyRepo    repository.HistoryRepository
}

func NewProductService(barcodeChecker BarcodeChecker, historyRepo repository.HistoryRepository) ProductService {
	return &productService{
		barcodeChecker: barcodeChecker,
		historyRepo:    historyRepo,
	}
}

func (s *productService) CheckProduct(barcode string, userID uint, clientIP string) (*model.ProductCheckResult, error) {
	isOriginal, err := s.barcodeChecker.Check(barcode)
	if err != nil {
		return nil, err
	}

	barcodeCountry := utils.GetCountryFromBarcode(barcode)
	var ipCountry string
	if clientIP != "" {
		country, err := utils.GetCountryFromIP(clientIP)
		if err == nil {
			ipCountry = country
		}
	}

	result := &model.ProductCheckResult{
		Barcode:        barcode,
		IsOriginal:     isOriginal,
		BarcodeCountry: barcodeCountry,
		IPCountry:      ipCountry,
		CheckedAt:      time.Now(),
	}

	if userID != 0 {
		check := &model.ProductCheck{
			UserID:         userID,
			Barcode:        barcode,
			IsOriginal:     isOriginal,
			BarcodeCountry: barcodeCountry,
			IPCountry:      ipCountry,
			CheckedAt:      time.Now(),
		}

		if err := s.historyRepo.Create(check); err != nil {
			return result, err
		}
		result.ID = check.ID.Hex()
	}

	return result, nil
}
