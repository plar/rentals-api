package mocks

import (
	"github.com/plar/rentals-api/domain"
	"github.com/plar/rentals-api/service"

	"github.com/stretchr/testify/mock"
)

type RentalService struct {
	mock.Mock
}

var _ service.RentalService = (*RentalService)(nil)

func (s *RentalService) GetAllRentals() ([]domain.Rental, error) {
	args := s.Called()
	return args.Get(0).([]domain.Rental), args.Error(1)
}

func (s *RentalService) GetRentalByID(id uint) (domain.Rental, error) {
	args := s.Called(id)
	return args.Get(0).(domain.Rental), args.Error(1)
}

func (s *RentalService) GetRentalsByFilter(filter domain.RentalFindFilter) (domain.Response[domain.Rental], error) {
	args := s.Called(filter)
	return args.Get(0).(domain.Response[domain.Rental]), args.Error(1)
}

// Add more methods as needed
