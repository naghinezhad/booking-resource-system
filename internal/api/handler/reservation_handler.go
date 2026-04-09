package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/naghinezhad/BookingResourceSystem/internal/service"
	"go.uber.org/zap"
)

type ReservationHandler struct {
	service *service.ReservationService
	logger  *zap.Logger
}

func NewReservationHandler(s *service.ReservationService, logger *zap.Logger) *ReservationHandler {
	return &ReservationHandler{service: s, logger: logger}
}

type ReserveRequest struct {
	ResourceID string `json:"resource_id" binding:"required"`
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
}

func (h *ReservationHandler) Reserve(c *gin.Context) {

	var req ReserveRequest

	if err := c.ShouldBindJSON(&req); err != nil {

		h.logger.Warn("invalid reserve request",
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("reserve request received",
		zap.String("resource_id", req.ResourceID),
		zap.String("start_time", req.StartTime),
		zap.String("end_time", req.EndTime),
	)

	start, err := time.Parse(time.RFC3339, req.StartTime)
	if err != nil {

		h.logger.Warn("invalid start_time format",
			zap.String("start_time", req.StartTime),
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid start_time"})
		return
	}

	end, err := time.Parse(time.RFC3339, req.EndTime)
	if err != nil {

		h.logger.Warn("invalid end_time format",
			zap.String("end_time", req.EndTime),
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid end_time"})
		return
	}

	err = h.service.Reserve(c, req.ResourceID, start, end)

	if err != nil {

		h.logger.Info("reservation conflict",
			zap.String("resource_id", req.ResourceID),
			zap.Time("start_time", start),
			zap.Time("end_time", end),
			zap.Error(err),
		)

		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("reservation created successfully",
		zap.String("resource_id", req.ResourceID),
		zap.Time("start_time", start),
		zap.Time("end_time", end),
	)

	c.JSON(http.StatusCreated, gin.H{"message": "reservation created"})
}

func (h *ReservationHandler) GetReservations(c *gin.Context) {

	id := c.Query("resource_id")

	h.logger.Info("get reservations request",
		zap.String("resource_id", id),
	)

	res, err := h.service.GetReservations(c, id)

	if err != nil {

		h.logger.Error("failed to fetch reservations",
			zap.String("resource_id", id),
			zap.Error(err),
		)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("reservations fetched successfully",
		zap.String("resource_id", id),
		zap.Int("count", len(res)),
	)

	c.JSON(http.StatusOK, res)
}
