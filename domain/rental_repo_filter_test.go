package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRentalFindFilterGetters(t *testing.T) {
	filter := &RentalFindFilter{
		priceMin:  ptr[uint](1000),
		priceMax:  ptr[uint](2000),
		limit:     ptr[uint](10),
		offset:    ptr[uint](0),
		rentalIDs: []int{1, 2, 3},
		lat:       ptr(12.9715987),
		long:      ptr(77.5945627),
		sort:      &SortPriceAsc,
	}

	priceMin, ok := filter.PriceMin()
	assert.True(t, ok)
	assert.Equal(t, uint(1000), priceMin)

	priceMax, ok := filter.PriceMax()
	assert.True(t, ok)
	assert.Equal(t, uint(2000), priceMax)

	limit, ok := filter.Limit()
	assert.True(t, ok)
	assert.Equal(t, uint(10), limit)

	offset, ok := filter.Offset()
	assert.True(t, ok)
	assert.Equal(t, uint(0), offset)

	rentalIDs, ok := filter.RentalIDs()
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3}, rentalIDs)

	coords, ok := filter.Coords()
	assert.True(t, ok)
	assert.Equal(t, 12.9715987, coords[0])
	assert.Equal(t, 77.5945627, coords[1])

	sort, ok := filter.Sort()
	assert.True(t, ok)
	assert.Equal(t, SortPriceAsc, sort)
}

func ptr[T any](v T) *T {
	return &v
}

func TestRentalFindFilterBuilder(t *testing.T) {
	assert := assert.New(t)

	filter, err := NewRentalFilterBuilder().
		WithPriceMin(1000).
		WithPriceMax(2000).
		WithLimit(10).
		WithOffset(0).
		WithRentalIDs([]int{1, 2, 3}).
		WithCoords([2]float64{12.9715987, 77.5945627}).
		WithSort(SortPriceAsc).
		Build()

	assert.NoError(err)
	assert.Equal(ptr[uint](1000), filter.priceMin)
	assert.Equal(ptr[uint](2000), filter.priceMax)
	assert.Equal(ptr[uint](10), filter.limit)
	assert.Equal(ptr[uint](0), filter.offset)
	assert.Equal([]int{1, 2, 3}, filter.rentalIDs)
	assert.Equal(ptr(12.9715987), filter.lat)
	assert.Equal(ptr(77.5945627), filter.long)
	assert.Equal(ptr(SortPriceAsc), filter.sort)
}

func TestRentalFindFilterBuilderError(t *testing.T) {
	assert := assert.New(t)

	_, err := NewRentalFilterBuilder().
		WithPriceMin(2000).
		WithPriceMax(1000).
		Build()
	assert.Error(err)
}
