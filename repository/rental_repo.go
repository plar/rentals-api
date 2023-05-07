package repository

import (
	"github.com/plar/rentals-api/domain"
	"go.uber.org/zap"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rentalRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

var _ domain.RentalRepository = (*rentalRepository)(nil)

func NewRentalRepository(db *gorm.DB, logger *zap.Logger) domain.RentalRepository {
	if logger == nil {
		logger = zap.NewNop()
	}

	return &rentalRepository{
		db:     db,
		logger: logger,
	}
}

func RentalRepositoryMigrate(db *gorm.DB) {
	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Rental{})
}

func (r *rentalRepository) FindAll() ([]domain.Rental, error) {
	var rentals []Rental
	err := r.db.Preload("User").Find(&rentals).Error
	return toDomainRentals(rentals), err
}

func (r *rentalRepository) FindByID(id uint) (domain.Rental, error) {
	var rental Rental
	err := r.db.Preload("User").First(&rental, id).Error
	return toDomainRental(rental), err
}

func (r *rentalRepository) applyViewFilter(filter domain.ViewFilter) func(db *gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		if limit, ok := filter.Limit(); ok {
			query = query.Limit(int(limit))
		}
		if offset, ok := filter.Offset(); ok {
			query = query.Offset(int(offset))
		}
		if sort, sortOk := filter.Sort(); sortOk {
			query = query.Order(clause.OrderByColumn{Column: clause.Column{Name: sort.Field()}, Desc: sort.IsDesc()})
		}
		return query
	}
}

func (r *rentalRepository) applySelectionFilter(filter domain.RentalFindFilter) func(db *gorm.DB) *gorm.DB {
	return func(query *gorm.DB) *gorm.DB {
		// price
		priceMin, priceMinOk := filter.PriceMin()
		priceMax, priceMaxOk := filter.PriceMax()
		if priceMinOk && priceMaxOk {
			query = query.Where("price_per_day >= ? AND price_per_day <= ?", priceMin, priceMax)
		} else if priceMinOk {
			query = query.Where("price_per_day >= ?", priceMin)
		} else if priceMaxOk {
			query = query.Where("price_per_day <= ?", priceMax)
		}

		// ids
		if ids, ok := filter.RentalIDs(); ok {
			query = query.Where("id IN (?)", ids)
		}

		// near
		if near, nearOk := filter.Coords(); nearOk {
			// Calculate the distance between two points using the Haversine formula
			// You can adjust the distance value (100 miles) according to your needs
			query = query.Where("earth_distance(ll_to_earth(lat, lng), ll_to_earth(?, ?)) <= ?", near[0], near[1], 100*1609.34)
		}
		return query
	}
}

func (r *rentalRepository) FindByFilter(filter domain.RentalFindFilter) (domain.Response[domain.Rental], error) {
	// preload Users
	query := r.db.Preload("User")
	query = query.Scopes(r.applySelectionFilter(filter))

	// count total filtered items
	var (
		items []Rental
		total int64
	)
	query.Model(items).Count(&total)
	query = query.Scopes(r.applyViewFilter(&filter))
	// query filtered items
	err := query.Find(&items).Error

	return domain.NewResponse(&filter, total, items, toDomainRentals), err
}

func toDomainRental(r Rental) domain.Rental {
	dr := domain.Rental{
		ID:              r.ID,
		Name:            r.Name,
		Description:     r.Description,
		Type:            r.Type,
		Make:            r.Make,
		Model:           r.Model,
		Year:            r.Year,
		Length:          r.Length,
		Sleeps:          r.Sleeps,
		PrimaryImageURL: r.PrimaryImageURL,
		Price: domain.Price{
			Day: r.Price,
		},
		Location: domain.Location{
			City:    r.City,
			State:   r.State,
			Zip:     r.Zip,
			Country: r.Country,
			Lat:     r.Lat,
			Lng:     r.Lng,
		},
		User: domain.User{
			ID:        int(r.UserID),
			FirstName: r.User.FirstName,
			LastName:  r.User.LastName,
		},
	}
	return dr
}

func toDomainRentals(rs []Rental) (drs []domain.Rental) {
	if len(rs) == 0 {
		return []domain.Rental{}
	}
	for _, r := range rs {
		drs = append(drs, toDomainRental(r))
	}
	return
}
