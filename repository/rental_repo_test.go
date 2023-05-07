package repository_test

import (
	"database/sql"
	"log"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/plar/rentals-api/domain"
	"github.com/plar/rentals-api/repository"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/suite"
)

type RentalRepoTestSuite struct {
	suite.Suite

	godb   *sql.DB
	gormdb *gorm.DB
	mock   sqlmock.Sqlmock

	rental domain.Rental
	user   domain.User
}

func GormLogger() logger.Interface {
	logWriter := log.New(os.Stdout, "\r\n", log.LstdFlags)
	logCfg := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Info,
		IgnoreRecordNotFoundError: false,
		Colorful:                  false,
	}
	return logger.New(logWriter, logCfg)
}

func (s *RentalRepoTestSuite) BeforeTest(suiteName, testName string) {
	var (
		err error
	)
	s.godb, s.mock, err = sqlmock.New()
	s.Assertions.NoError(err, "Failed to open mock sql db")
	s.Assertions.NotNil(s.godb, "mock db is null")
	s.Assertions.NotNil(s.mock, "sqlmock is null")

	cfg := postgres.Config{
		DriverName: "postgres",
		Conn:       s.godb,
	}

	if s.gormdb, err = gorm.Open(postgres.New(cfg), &gorm.Config{
		Logger: GormLogger(),
	}); err != nil {
		s.Assertions.Fail("cannot open mock postgres DB")
	}
}

func (s *RentalRepoTestSuite) AfterTest(suiteName, testName string) {
	if s.godb != nil {
		s.godb.Close()
	}
}

func (s *RentalRepoTestSuite) TestFindAll() {
	format := "2006-01-02 15:04:05.999999-07"
	value := "2021-11-29 22:42:06.478595+00"
	createdAt, err := time.Parse(format, value)
	s.Assertions.NoError(err)

	// setup mock
	rentalRows := sqlmock.NewRows([]string{
		"id", "created", "updated", "user_id", "name", "type", "description", "sleeps", "price_per_day", "home_city", "home_state", "home_zip", "home_country", "vehicle_make", "vehicle_model", "vehicle_year", "vehicle_length", "lat", "lng", "primary_image_url",
	}).AddRow(1, createdAt, createdAt, 1,
		"'Abaco' VW Bay Window: Westfalia Pop-top",
		"camper-van",
		"ultrices consectetur torquent",
		4, 16900, "Costa Mesa", "CA", "92627", "US", "Volkswagen", "Bay Window", 1978, 15, 33.64, -117.93, "https://res.cloudinary.com/outdoorsy/image/upload/v1528586451/p/rentals/4447/images/yd7txtw4hnkjvklg8edg.jpg")

	s.mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT * 
	  FROM "rentals"
	`)).WillReturnRows(rentalRows)

	userRows := sqlmock.NewRows([]string{
		"id", "first_name", "last_name",
	}).AddRow(1, "John", "Smith")
	s.mock.ExpectQuery(regexp.QuoteMeta(`
	SELECT * 
	  FROM "users" 
	 WHERE "users"."id" = $1 AND "users"."deleted_at" IS NULL
	`)).WithArgs(1).WillReturnRows(userRows)

	// run repo test
	rentalRepo := repository.NewRentalRepository(s.gormdb, nil)
	actualRentals, err := rentalRepo.FindAll()

	// check asserts
	s.Assertions.NoError(err)
	s.Assertions.Len(actualRentals, 1)

	expectedRentals := []domain.Rental{
		{
			ID:          1,
			Name:        "'Abaco' VW Bay Window: Westfalia Pop-top",
			Description: "ultrices consectetur torquent",
			Type:        "camper-van",
			Make:        "Volkswagen",
			Model:       "Bay Window",
			Year:        1978,
			Length:      15,
			Sleeps:      4,
			Price: domain.Price{
				Day: 16900,
			},
			Location: domain.Location{
				City:    "Costa Mesa",
				State:   "CA",
				Zip:     "92627",
				Country: "US",
				Lat:     33.64,
				Lng:     -117.93,
			},
			PrimaryImageURL: "https://res.cloudinary.com/outdoorsy/image/upload/v1528586451/p/rentals/4447/images/yd7txtw4hnkjvklg8edg.jpg",
			User: domain.User{
				ID:        1,
				FirstName: "John",
				LastName:  "Smith",
			},
		},
	}

	s.Assertions.Equal(expectedRentals, actualRentals)

	s.Assertions.NoError(s.mock.ExpectationsWereMet(), "Failed to meet expectations")
}

func TestRentalRepoSuite(t *testing.T) {
	suite.Run(t, &RentalRepoTestSuite{})
}

// Add more repository layer tests
