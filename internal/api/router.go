package api

import (
	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/api/handler"
	"github.com/naghinezhad/BookingResourceSystem/internal/api/middleware"
	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/concurrency"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
	"github.com/naghinezhad/BookingResourceSystem/internal/lock"
	"github.com/naghinezhad/BookingResourceSystem/internal/logger"
	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(
	mongo *database.Mongo,
	redis *cache.Redis,
	cfg *config.Config,
) *gin.Engine {
	metrics.Register()

	r := gin.Default()

	// metrics
	r.Use(middleware.MetricsMiddleware())

	// logger
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

	// Request Limiter
	limiter := concurrency.NewRequestLimiter(cfg.MaxConcurrentRequests)
	r.Use(middleware.RequestLimitMiddleware(limiter))

	// Worker Pool
	workerPool := concurrency.NewWorkerPool(
		3,
		cfg.MaxConcurrentRequests,
		reservationService,
	)

	// handlers
	reservationHandler := handler.NewReservationHandler(
		reservationService,
		workerPool,
		log,
	)

	availabilityHandler := handler.NewAvailabilityHandler(
		availabilityService,
		log,
	)

	// routes
	r.POST("/reserve", reservationHandler.Reserve)
	r.GET("/availability", availabilityHandler.Check)
	r.GET("/reservations", reservationHandler.GetReservations)
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	return r
}
