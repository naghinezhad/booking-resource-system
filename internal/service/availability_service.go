package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AvailabilityService struct {
	repo  *repository.ReservationRepository
	cache *cache.Redis
}

func NewAvailabilityService(
	repo *repository.ReservationRepository,
	cache *cache.Redis,
) *AvailabilityService {

	return &AvailabilityService{
		repo:  repo,
		cache: cache,
	}
}

func (s *AvailabilityService) CheckAvailability(
	ctx context.Context,
	id string,
	start time.Time,
	end time.Time,
) (bool, error) {

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, errors.New("invalid resource id")
	}

	cacheKey := fmt.Sprintf("availability:%s", id)

	// Try Redis cache first
	val, err := s.cache.Client.Get(ctx, cacheKey).Result()
	if err == nil {
		return val == "1", nil
	}

	available, err := s.repo.CheckAvailability(ctx, objID, start, end)
	if err != nil {
		return false, err
	}

	// cache it
	s.cache.Client.Set(ctx, cacheKey, map[bool]string{true: "1", false: "0"}[available], 5*time.Second)

	return available, nil
}
