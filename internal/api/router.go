package api

import (
	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/api/handler"
	"github.com/naghinezhad/BookingResourceSystem/internal/api/middleware"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
	"github.com/naghinezhad/BookingResourceSystem/internal/redis"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
	"github.com/naghinezhad/BookingResourceSystem/internal/utils/concurrency"
	"github.com/naghinezhad/BookingResourceSystem/internal/utils/logger"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func SetupRouter(
	mongo *database.Mongo,
	redisClient redis.Client,
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

	// services
	reservationService := service.NewReservationService(
		reservationRepo,
		redisClient,
	)

	availabilityService := service.NewAvailabilityService(
		reservationRepo,
		redisClient,
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
