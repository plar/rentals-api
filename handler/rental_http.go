package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/plar/rentals-api/domain"
	"github.com/plar/rentals-api/service"
)

type RentalHandler interface {
	GetRentalByID(c *gin.Context)
	GetRentals(c *gin.Context)
}

type rentalHandler struct {
	service service.RentalService
	logger  *zap.Logger
}

func NewRentalHandler(service service.RentalService, logger *zap.Logger) RentalHandler {
	if logger == nil {
		logger = zap.NewNop()
	}
	return &rentalHandler{
		service: service,
		logger:  logger,
	}
}

type RentalByIDRequest struct {
	ID uint `uri:"id"`
}

func (h *rentalHandler) GetRentalByID(c *gin.Context) {
	var req RentalByIDRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid rental ID"})
		return
	}
	rental, err := h.service.GetRentalByID(req.ID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "rental not found"})
		return
	}

	c.JSON(http.StatusOK, rental)
}

type RentalsRequest struct {
	PriceMin *uint   `form:"price_min" binding:"omitempty,gte=0"`
	PriceMax *uint   `form:"price_max" binding:"omitempty,gte=0"`
	Limit    *uint   `form:"limit,default=10" binding:"min=1,max=100"`
	Offset   *uint   `form:"offset,default=0" binding:"omitempty,gte=0"`
	IDs      *string `form:"ids"`
	Near     *string `form:"near"`
	Sort     *string `form:"sort" binding:"omitempty,oneof=price price_asc price_desc year year_asc year_desc"`

	parsedIDs  []int
	parsedNear [2]float64
	parsedSort domain.Sort
}

func (r *RentalsRequest) validate() (err error) {
	if r.IDs != nil {
		if r.parsedIDs, err = toIntSlice(*r.IDs); err != nil {
			return fmt.Errorf("invalid ids input: %w", err)
		}
	}

	if r.Near != nil {
		if r.parsedNear, err = toNear(*r.Near); err != nil {
			return fmt.Errorf("invalid near input: %w", err)
		}
	}

	r.parsedSort = domain.SortNone
	if r.Sort != nil {
		r.parsedSort = map[string]domain.Sort{
			"price":      domain.SortPriceAsc,
			"price_asc":  domain.SortPriceAsc,
			"price_desc": domain.SortPriceDesc,
			"year":       domain.SortYearAsc,
			"year_asc":   domain.SortYearAsc,
			"year_desc":  domain.SortYearDesc,
		}[*r.Sort]
	}

	return nil
}

func createRentalFindFilter(c *gin.Context) (filter domain.RentalFindFilter, err error) {
	var req RentalsRequest
	if err = c.ShouldBind(&req); err != nil {
		return
	} else if err = req.validate(); err != nil {
		return
	}

	return toDomainRentalFilter(req)
}

func (h *rentalHandler) GetRentals(c *gin.Context) {
	filter, err := createRentalFindFilter(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.GetRentalsByFilter(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func toIntSlice(s string) (ints []int, _ error) {
	strs := strings.Split(s, ",")
	for _, s := range strs {
		n, nerr := strconv.ParseInt(strings.TrimSpace(s), 10, 64)
		if nerr != nil {
			return nil, fmt.Errorf("error parsing number %q: %w", s, nerr)
		}
		ints = append(ints, int(n))
	}
	return
}

func toNear(s string) (near [2]float64, _ error) {
	strs := strings.Split(s, ",")
	if len(strs) != 2 {
		return [2]float64{}, fmt.Errorf("invalid number of coords")
	}
	for i, s := range strs {
		n, nerr := strconv.ParseFloat(strings.TrimSpace(s), 64)
		if nerr != nil {
			return [2]float64{}, fmt.Errorf("error parsing coord %q: %w", s, nerr)
		}
		near[i] = float64(n)
	}
	return
}

func toDomainRentalFilter(inp RentalsRequest) (domain.RentalFindFilter, error) {
	b := domain.NewRentalFilterBuilder()
	if inp.PriceMin != nil {
		b.WithPriceMin(*inp.PriceMin)
	}

	if inp.PriceMax != nil {
		b.WithPriceMax(*inp.PriceMax)
	}

	if inp.Limit != nil {
		b.WithLimit(*inp.Limit)
	}

	if inp.Offset != nil {
		b.WithOffset(*inp.Offset)
	}

	if inp.IDs != nil && len(*inp.IDs) > 0 {
		b.WithRentalIDs(inp.parsedIDs)
	}

	if inp.Near != nil {
		b.WithCoords(inp.parsedNear)
	}

	if inp.Sort != nil {
		b.WithSort(inp.parsedSort)
	}

	return b.Build()
}
