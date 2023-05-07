package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/zap"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"moul.io/zapgorm2"

	"github.com/plar/rentals-api/config"
	"github.com/plar/rentals-api/handler"
	"github.com/plar/rentals-api/logs"
	"github.com/plar/rentals-api/repository"
	"github.com/plar/rentals-api/service"
)

func main() {
	log := logs.Init()
	defer log.Sync()

	// configure gorm logger to use zap logger
	gormLogger := zapgorm2.New(log)
	gormLogger.SetAsDefault()
	// ... and create gorm
	db, err := gorm.Open(postgres.Open(config.DBConnectionString()), &gorm.Config{
		Logger: gormLogger.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database", zap.Error(err))
	}

	// close DB connection
	defer func() {
		dbInst, err := db.DB()
		if err != nil {
			log.Error("Cannot get DB", zap.Error(err))
		}
		if err = dbInst.Close(); err != nil {
			log.Error("Cannot close DB", zap.Error(err))
		}
	}()

	// setup app
	rentalRepo := repository.NewRentalRepository(db, log)
	rentalRepoLog := repository.NewRentalRepositoryLogger(rentalRepo, log)
	rentalSvc := service.NewRentalService(rentalRepoLog, log)
	rentalHandler := handler.NewRentalHandler(rentalSvc, log)

	// run migrations
	repository.RentalRepositoryMigrate(db)

	router := gin.New()
	router.Use(ginzap.Ginzap(log, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(log, true))

	// regisyer handlers
	router.GET("/rentals/:id", rentalHandler.GetRentalByID)
	router.GET("/rentals", rentalHandler.GetRentals)

	// run HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// handle SYS signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// wait for quit signal...
	<-quit

	// enforce shutdown in 5s
	log.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", zap.Error(err))
	}

	log.Info("Server exiting")
}
