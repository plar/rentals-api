package repository

import (
	"fmt"

	"github.com/plar/rentals-api/domain"
	"go.uber.org/zap"
)

type rentalRepositoryLogger struct {
	next   domain.RentalRepository
	logger *zap.Logger
}

var _ domain.RentalRepository = (*rentalRepositoryLogger)(nil)

func NewRentalRepositoryLogger(next domain.RentalRepository, logger *zap.Logger) domain.RentalRepository {
	return &rentalRepositoryLogger{
		next:   next,
		logger: logger,
	}
}

func (l *rentalRepositoryLogger) FindAll() (rentals []domain.Rental, err error) {
	l.logger.Debug("FindAll called")
	defer func() {
		if err == nil {
			l.logger.Debug("FindAll completed")
		} else {
			l.logger.Error("FindAll error", zap.Error(err))
		}
	}()
	return l.next.FindAll()
}

func (l *rentalRepositoryLogger) FindByID(id uint) (rental domain.Rental, err error) {
	l.logger.Debug("FindByID called", zap.Uint("id", id))
	defer func() {
		if err == nil {
			l.logger.Debug("FindByID completed", zap.String("rental", fmt.Sprintf("%v", rental)))
		} else {
			l.logger.Error("FindByID error", zap.Error(err))
		}
	}()
	return l.next.FindByID(id)
}

func (l *rentalRepositoryLogger) FindByFilter(filter domain.RentalFindFilter) (rentals domain.Response[domain.Rental], err error) {
	l.logger.Debug("FindByFilter called", zap.String("filter", filter.String()))
	defer func() {
		if err == nil {
			l.logger.Debug("FindByFilter completed")
		} else {
			l.logger.Error("FindByFilter error", zap.Error(err))
		}
	}()
	return l.next.FindByFilter(filter)
}

// Add more methods as needed
