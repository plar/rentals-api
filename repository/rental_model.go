package repository

import (
	"time"

	"gorm.io/gorm"
)

type Rental struct {
	ID        uint           `gorm:"primary_key"`
	CreatedAt time.Time      `gorm:"column:created;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	UserID uint `gorm:"column:user_id"`
	User   User

	Name            string
	Description     string
	Type            string
	Make            string  `gorm:"column:vehicle_make"`
	Model           string  `gorm:"column:vehicle_model"`
	Year            int     `gorm:"column:vehicle_year"`
	Length          float64 `gorm:"column:vehicle_length"`
	Sleeps          int
	Price           int    `gorm:"column:price_per_day"`
	City            string `gorm:"column:home_city"`
	State           string `gorm:"column:home_state"`
	Zip             string `gorm:"column:home_zip"`
	Country         string `gorm:"column:home_country"`
	PrimaryImageURL string `gorm:"column:primary_image_url"`
	Lat             float64
	Lng             float64
}

type User struct {
	gorm.Model

	FirstName string
	LastName  string
}

func (r *Rental) TableName() string {
	return "rentals"
}
