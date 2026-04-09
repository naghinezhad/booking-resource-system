package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
)

type AvailabilityHandler struct {
	service *service.AvailabilityService
}

func NewAvailabilityHandler(s *service.AvailabilityService) *AvailabilityHandler {
	return &AvailabilityHandler{service: s}
}

func (h *AvailabilityHandler) Check(c *gin.Context) {

	resourceID := c.Query("resource_id")
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	available, err := h.service.CheckAvailability(
		c.Request.Context(),
		resourceID,
		start,
		end,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": available,
	})
}
