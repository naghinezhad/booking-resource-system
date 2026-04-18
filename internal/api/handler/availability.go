package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
	"go.uber.org/zap"
)

type AvailabilityHandler struct {
	service *service.AvailabilityService
	logger  *zap.Logger
}

func NewAvailabilityHandler(s *service.AvailabilityService, logger *zap.Logger) *AvailabilityHandler {
	return &AvailabilityHandler{service: s, logger: logger}
}

func (h *AvailabilityHandler) Check(c *gin.Context) {

	resourceID := c.Query("resource_id")
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

	h.logger.Info("availability check requested",
		zap.String("resource_id", resourceID),
		zap.String("start_time", startStr),
		zap.String("end_time", endStr),
	)

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		h.logger.Warn("invalid start_time format",
			zap.String("start_time", startStr),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		h.logger.Warn("invalid end_time format",
			zap.String("end_time", endStr),
			zap.Error(err),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	available, err := h.service.CheckAvailability(c, resourceID, start, end)
	if err != nil {
		h.logger.Error("availability check failed",
			zap.String("resource_id", resourceID),
			zap.Error(err),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("availability check success",
		zap.String("resource_id", resourceID),
		zap.Bool("available", available),
	)

	c.JSON(http.StatusOK, gin.H{"available": available})
}
