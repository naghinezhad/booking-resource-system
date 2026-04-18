package service

import (
	"context"
	"fmt"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
)

type AvailabilityService struct {
	repo  *repository.ReservationRepository
	cache cache.Client
}

func NewAvailabilityService(
	repo *repository.ReservationRepository,
	cacheClient cache.Client,
) *AvailabilityService {

	return &AvailabilityService{
		repo:  repo,
		cache: cacheClient,
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

	val, err := s.cache.Get(ctx, cacheKey)
	if err == nil {

		metrics.CacheHits.Inc()

		return val == "true", nil
	}

	metrics.CacheMiss.Inc()

	available, err := s.repo.CheckAvailability(ctx, id, start, end)
	if err != nil {
		return false, err
	}

	_ = s.cache.Set(ctx, cacheKey, fmt.Sprintf("%t", available), 5*time.Second)

	return available, nil
}
