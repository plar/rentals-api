package service_test

import (
	"testing"

	"github.com/plar/rentals-api/domain"
	"github.com/plar/rentals-api/repository/mocks"
	"github.com/plar/rentals-api/service"

	"github.com/stretchr/testify/assert"
)

func TestGetAllRentals(t *testing.T) {
	mockRepo := &mocks.RentalRepository{}
	mockRentals := []domain.Rental{
		{ID: 1, Name: "Test Rental 1"},
		{ID: 2, Name: "Test Rental 2"},
	}

	mockRepo.On("FindAll").Return(mockRentals, nil)

	rentalService := service.NewRentalService(mockRepo, nil)
	rentals, err := rentalService.GetAllRentals()

	assert.NoError(t, err)
	assert.Len(t, rentals, 2)
	mockRepo.AssertExpectations(t)
}

// Add more service layer tests
