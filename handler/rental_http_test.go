package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/plar/rentals-api/domain"
	"github.com/plar/rentals-api/handler"
	"github.com/plar/rentals-api/service/mocks"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetRentalByID(t *testing.T) {
	mockRental := domain.Rental{ID: 1, Name: "Test Rental"}

	t.Run("/rentals/1", func(t *testing.T) {
		mockService := new(mocks.RentalService)
		mockService.On("GetRentalByID", mock.AnythingOfType("uint")).Return(mockRental, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := handler.NewRentalHandler(mockService, nil)
		router.GET("/rentals/:id", handler.GetRentalByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rentals/1", bytes.NewBuffer(nil))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var retrievedRental domain.Rental
		err := json.Unmarshal(w.Body.Bytes(), &retrievedRental)
		assert.NoError(t, err)

		assert.Equal(t, mockRental.ID, retrievedRental.ID)
		assert.Equal(t, mockRental.Name, retrievedRental.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("/rentals/badid", func(t *testing.T) {
		mockService := new(mocks.RentalService)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := handler.NewRentalHandler(mockService, nil)
		router.GET("/rentals/:id", handler.GetRentalByID)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rentals/badid", bytes.NewBuffer(nil))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)

		var errResp struct {
			Error string
		}
		err := json.Unmarshal(w.Body.Bytes(), &errResp)
		assert.NoError(t, err)
		assert.Equal(t, "invalid rental ID", errResp.Error)

		mockService.AssertExpectations(t)
	})

}

func TestGetRentals(t *testing.T) {
	mockResponse := domain.Response[domain.Rental]{
		Paginator: domain.Paginator{
			Limit:      0,
			Offset:     0,
			TotalItems: 2,
		},
		Items: []domain.Rental{
			{ID: 1, Name: "Test Rental1"},
			{ID: 2, Name: "Test Rental2"},
		},
	}

	t.Run("/rentals (check filter default values)", func(t *testing.T) {
		mockService := &mocks.RentalService{}
		mockService.On("GetRentalsByFilter", mock.MatchedBy(func(f domain.RentalFindFilter) bool {
			// TODO: write a func to compare RentalFindFilter
			if _, ok := f.PriceMin(); ok {
				return false
			}
			if _, ok := f.PriceMax(); ok {
				return false
			}
			if limit, ok := f.Limit(); !ok || limit != uint(10) {
				return false
			}
			if offset, ok := f.Offset(); !ok || offset != uint(0) {
				return false
			}
			if _, ok := f.RentalIDs(); ok {
				return false
			}
			if _, ok := f.Coords(); ok {
				return false
			}
			if _, ok := f.Sort(); ok {
				return false
			}
			return true
		})).Return(mockResponse, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := handler.NewRentalHandler(mockService, nil)
		router.GET("/rentals", handler.GetRentals)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rentals", bytes.NewBuffer(nil))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse domain.Response[domain.Rental]
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, actualResponse)

		mockService.AssertExpectations(t)
	})

	t.Run("/rentals (check all filter values)", func(t *testing.T) {
		mockService := &mocks.RentalService{}
		mockService.On("GetRentalsByFilter", mock.MatchedBy(func(f domain.RentalFindFilter) bool {
			// TODO: write a func to compare RentalFindFilter
			if pmin, ok := f.PriceMin(); !ok || pmin != uint(100) {
				return false
			}
			if pmax, ok := f.PriceMax(); !ok || pmax != uint(200) {
				return false
			}
			if limit, ok := f.Limit(); !ok || limit != uint(11) {
				return false
			}
			if offset, ok := f.Offset(); !ok || offset != uint(22) {
				return false
			}
			if ids, ok := f.RentalIDs(); !ok || len(ids) != 3 || ids[0] != 1 || ids[1] != 2 || ids[2] != 3 {
				return false
			}
			if coords, ok := f.Coords(); !ok || coords[0] != 1.23 || coords[1] != 4.56 {
				return false
			}
			if sort, ok := f.Sort(); !ok || sort != domain.SortPriceAsc {
				return false
			}
			return true
		})).Return(mockResponse, nil)

		gin.SetMode(gin.TestMode)
		router := gin.Default()
		handler := handler.NewRentalHandler(mockService, nil)
		router.GET("/rentals", handler.GetRentals)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/rentals?price_min=100&price_max=200&limit=11&offset=22&ids=1,2,3&near=1.23,4.56&sort=price_asc", bytes.NewBuffer(nil))
		router.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)

		var actualResponse domain.Response[domain.Rental]
		err := json.Unmarshal(w.Body.Bytes(), &actualResponse)
		assert.NoError(t, err)
		assert.Equal(t, mockResponse, actualResponse)

		mockService.AssertExpectations(t)
	})

}

// Add more handler layer tests
