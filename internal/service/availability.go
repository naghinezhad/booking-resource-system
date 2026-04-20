package service

import (
	"context"
	"fmt"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
	"github.com/naghinezhad/BookingResourceSystem/internal/redis"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
)

type AvailabilityService struct {
	repo  *repository.ReservationRepository
	redis redis.Client
}

func NewAvailabilityService(
	repo *repository.ReservationRepository,
	redisClient redis.Client,
) *AvailabilityService {

	return &AvailabilityService{
		repo:  repo,
		redis: redisClient,
	}
}

func (s *AvailabilityService) CheckAvailability(
	ctx context.Context,
	id string,
	start time.Time,
	end time.Time,
) (bool, error) {
	cacheKey := fmt.Sprintf(
		"availability:%s:%d:%d",
		id,
		start.Unix(),
		end.Unix(),
	)

	val, err := s.redis.Get(ctx, cacheKey)
	if err == nil {

		metrics.CacheHits.Inc()

		return val == "true", nil
	}

	metrics.CacheMiss.Inc()

	available, err := s.repo.CheckAvailability(ctx, id, start, end)
	if err != nil {
		return false, err
	}

	_ = s.redis.Set(ctx, cacheKey, fmt.Sprintf("%t", available), 5*time.Second)

	return available, nil
}
