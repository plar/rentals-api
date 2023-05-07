package service

import (
	"github.com/plar/rentals-api/domain"
	"go.uber.org/zap"
)

type RentalService interface {
	GetAllRentals() ([]domain.Rental, error)
	GetRentalByID(id uint) (domain.Rental, error)
	GetRentalsByFilter(filter domain.RentalFindFilter) (domain.Response[domain.Rental], error)
}

type rentalService struct {
	repo domain.RentalRepository
}

func NewRentalService(repo domain.RentalRepository, logger *zap.Logger) RentalService {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &rentalService{repo}
}

func (s *rentalService) GetAllRentals() ([]domain.Rental, error) {
	return s.repo.FindAll()
}

func (s *rentalService) GetRentalByID(id uint) (domain.Rental, error) {
	return s.repo.FindByID(id)
}

func (s *rentalService) GetRentalsByFilter(filter domain.RentalFindFilter) (domain.Response[domain.Rental], error) {
	return s.repo.FindByFilter(filter)
}
