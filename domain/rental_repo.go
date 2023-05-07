package domain

type RentalRepository interface {
	FindAll() ([]Rental, error)
	FindByID(id uint) (Rental, error)
	FindByFilter(filter RentalFindFilter) (Response[Rental], error)
}
