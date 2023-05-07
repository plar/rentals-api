package domain

import (
	"errors"
	"fmt"
	"strings"
)

type LimitFilter interface {
	Limit() (uint, bool)
}

type OffsetFilter interface {
	Offset() (uint, bool)
}

type SortFilter interface {
	Sort() (Sort, bool)
}

type ViewFilter interface {
	LimitFilter
	OffsetFilter
	SortFilter
}

// enum, use struct instead of type Sort string to avoid Sort("any sort type")
type Sort struct {
	s string
}

func (s Sort) Field() string {
	return strings.Split(s.s, "#")[0]
}

func (s Sort) Order() string {
	return strings.Split(s.s, "#")[1]
}

func (s Sort) IsDesc() bool {
	return strings.Split(s.s, "#")[1] == "desc"
}

func (s Sort) IsZero() bool {
	return len(s.s) == 0
}

var (
	SortNone      = Sort{}
	SortPriceAsc  = Sort{"price_per_day#asc"}
	SortPriceDesc = Sort{"price_per_day#desc"}
	SortYearAsc   = Sort{"vehicle_year#asc"}
	SortYearDesc  = Sort{"vehicle_year#desc"}
)

type RentalFindFilter struct {
	priceMin  *uint
	priceMax  *uint
	limit     *uint
	offset    *uint
	rentalIDs []int
	lat       *float64
	long      *float64
	sort      *Sort
}

var _ ViewFilter = (*RentalFindFilter)(nil)

func (b *RentalFindFilter) validate() error {
	pmin, pminOk := b.PriceMin()
	pmax, pmaxOk := b.PriceMax()
	if pminOk && pmaxOk && pmin > pmax {
		return errors.New("invalid priceMin and priceMax: priceMin > priceMax")
	}
	return nil
}

func (f *RentalFindFilter) PriceMin() (uint, bool) {
	if f.priceMin != nil {
		return *f.priceMin, true
	}
	return 0, false
}

func (f *RentalFindFilter) PriceMax() (uint, bool) {
	if f.priceMax != nil {
		return *f.priceMax, true
	}
	return 0, false
}

func (f *RentalFindFilter) Limit() (uint, bool) {
	if f.limit != nil {
		return *f.limit, true
	}
	return 0, false
}

func (f *RentalFindFilter) Offset() (uint, bool) {
	if f.offset != nil {
		return *f.offset, true
	}
	return 0, false
}

func (f *RentalFindFilter) RentalIDs() ([]int, bool) {
	if len(f.rentalIDs) > 0 {
		return f.rentalIDs, true
	}
	return nil, false
}

func (f *RentalFindFilter) Coords() ([2]float64, bool) {
	if f.lat != nil && f.long != nil {
		return [2]float64{*f.lat, *f.long}, true
	}
	return [2]float64{}, false
}

func (f *RentalFindFilter) Sort() (Sort, bool) {
	if f.sort != nil {
		return *f.sort, true
	}
	return SortNone, false
}

func (f *RentalFindFilter) String() string {
	kv := make(map[string]any)

	// price
	if priceMin, priceMinOk := f.PriceMin(); priceMinOk {
		kv["priceMin"] = priceMin
	}

	if priceMax, priceMaxOk := f.PriceMax(); priceMaxOk {
		kv["priceMax"] = priceMax
	}

	// ids
	if ids, ok := f.RentalIDs(); ok {
		kv["ids"] = ids
	}

	if limit, ok := f.Limit(); ok {
		kv["limit"] = limit
	}
	if offset, ok := f.Offset(); ok {
		kv["offset"] = offset
	}

	// near
	if near, nearOk := f.Coords(); nearOk {
		kv["near"] = near
	}

	// sort
	if sort, sortOk := f.Sort(); sortOk {
		kv["sort"] = sort
	}

	s := fmt.Sprintf("%v", kv)
	s = strings.TrimPrefix(s, "map[")
	s = strings.TrimSuffix(s, "]")
	return s
}

type RentalFindFilterBuilder struct {
	filter RentalFindFilter
}

func NewRentalFilterBuilder() *RentalFindFilterBuilder {
	return &RentalFindFilterBuilder{
		filter: RentalFindFilter{},
	}
}

func (b *RentalFindFilterBuilder) WithPriceMin(priceMin uint) *RentalFindFilterBuilder {
	b.filter.priceMin = &priceMin
	return b
}

func (b *RentalFindFilterBuilder) WithPriceMax(priceMax uint) *RentalFindFilterBuilder {
	b.filter.priceMax = &priceMax
	return b
}

func (b *RentalFindFilterBuilder) WithLimit(limit uint) *RentalFindFilterBuilder {
	b.filter.limit = &limit
	return b
}

func (b *RentalFindFilterBuilder) WithOffset(offset uint) *RentalFindFilterBuilder {
	b.filter.offset = &offset
	return b
}

func (b *RentalFindFilterBuilder) WithRentalIDs(rentalIDs []int) *RentalFindFilterBuilder {
	b.filter.rentalIDs = rentalIDs
	return b
}

func (b *RentalFindFilterBuilder) WithCoords(coords [2]float64) *RentalFindFilterBuilder {
	b.filter.lat = &coords[0]
	b.filter.long = &coords[1]
	return b
}

func (b *RentalFindFilterBuilder) WithSort(sort Sort) *RentalFindFilterBuilder {
	b.filter.sort = &sort
	return b
}

func (b *RentalFindFilterBuilder) Build() (_ RentalFindFilter, err error) {
	if err = b.filter.validate(); err != nil {
		return
	}
	return b.filter, nil
}
