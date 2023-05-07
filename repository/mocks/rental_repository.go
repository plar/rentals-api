package mocks

import (
	"github.com/plar/rentals-api/domain"

	"github.com/stretchr/testify/mock"
)

type RentalRepository struct {
	mock.Mock
}

var _ domain.RentalRepository = (*RentalRepository)(nil)

func (r *RentalRepository) FindAll() ([]domain.Rental, error) {
	args := r.Called()
	return args.Get(0).([]domain.Rental), args.Error(1)
}

func (r *RentalRepository) FindByID(id uint) (domain.Rental, error) {
	args := r.Called(id)
	return args.Get(0).(domain.Rental), args.Error(1)
}

func (r *RentalRepository) FindByFilter(filter domain.RentalFindFilter) (domain.Response[domain.Rental], error) {
	args := r.Called(filter)
	return args.Get(0).(domain.Response[domain.Rental]), args.Error(1)
}

// Add more methods as needed
