package api

import (
	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/config"
	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/database"
)

func SetupRouter(
	mongo *database.Mongo,
	redis *cache.Redis,
	cfg *config.Config,
) *gin.Engine {
	r := gin.Default()

	return r
}
