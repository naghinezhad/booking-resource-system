package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/naghinezhad/BookingResourceSystem/internal/cache"
	"github.com/naghinezhad/BookingResourceSystem/internal/lock"
	"github.com/naghinezhad/BookingResourceSystem/internal/metrics"
	"github.com/naghinezhad/BookingResourceSystem/internal/model"
	"github.com/naghinezhad/BookingResourceSystem/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ReservationService struct {
	repo  *repository.ReservationRepository
	cache *cache.Redis
	lock  *lock.RedisLock
}

func NewReservationService(
	repo *repository.ReservationRepository,
	cache *cache.Redis,
	lock *lock.RedisLock,
) *ReservationService {
	return &ReservationService{
		repo:  repo,
		cache: cache,
		lock:  lock,
	}
}

func (s *ReservationService) Reserve(
	ctx context.Context,
	resourceID string,
	start time.Time,
	end time.Time,
) error {

	objID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("invalid").Inc()

		return errors.New("invalid resource id")
	}

	// distributed lock
	lockKey := fmt.Sprintf("lock:resource:%s", resourceID)
	acquired, err := s.lock.Acquire(ctx, lockKey, 3*time.Second)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	if !acquired {

		metrics.ReservationsTotal.WithLabelValues("busy").Inc()

		return errors.New("resource busy, try again")
	}

	defer func() {
		if err := s.lock.Release(ctx, lockKey); err != nil {
			fmt.Printf("failed to release lock: %v\n", err)
		}
	}()

	// check availability (db)
	available, err := s.repo.CheckAvailability(ctx, objID, start, end)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	if !available {

		metrics.ReservationsTotal.WithLabelValues("conflict").Inc()

		return errors.New("resource not available")
	}

	// create reservation
	reservation := &model.Reservation{
		ResourceID: objID,
		StartTime:  start,
		EndTime:    end,
		CreatedAt:  time.Now(),
	}

	err = s.repo.Create(ctx, reservation)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	// invalidate redis cache
	cacheKey := fmt.Sprintf("availability:%s", resourceID)
	s.cache.Client.Del(ctx, cacheKey)

	metrics.ReservationsTotal.WithLabelValues("success").Inc()

	return nil
}

func (s *ReservationService) GetReservations(
	ctx context.Context,
	resourceID string,
) ([]model.Reservation, error) {

	objID, err := primitive.ObjectIDFromHex(resourceID)
	if err != nil {
		return nil, errors.New("invalid resource id")
	}

	return s.repo.GetByResourceID(ctx, objID)
}
