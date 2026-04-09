package api

import (
	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/api/handler"
	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
	"github.com/naghinezhad/BookingResourceSystem/internal/lock"
	"github.com/naghinezhad/BookingResourceSystem/internal/logger"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
)

func SetupRouter(
	mongo *database.Mongo,
	redis *cache.Redis,
	cfg *config.Config,
) *gin.Engine {

	r := gin.Default()

	// load global logger
	log := logger.Log

	// repositories
	reservationRepo := repository.NewReservationRepository(mongo.DB)

	// locking
	distLock := lock.NewRedisLock(redis.Client)

	// services
	reservationService := service.NewReservationService(
		reservationRepo,
		redis,
		distLock,
	)

	availabilityService := service.NewAvailabilityService(
		reservationRepo,
		redis,
	)

	// handlers (now with logger)
	reservationHandler := handler.NewReservationHandler(reservationService, log)
	availabilityHandler := handler.NewAvailabilityHandler(availabilityService, log)

	// routes
	r.POST("/reserve", reservationHandler.Reserve)
	r.GET("/availability", availabilityHandler.Check)
	r.GET("/reservations", reservationHandler.GetReservations)

	return r
}
