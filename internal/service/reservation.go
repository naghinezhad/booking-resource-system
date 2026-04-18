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
)

type ReservationService struct {
	repo  *repository.ReservationRepository
	cache cache.Client
	lock  lock.Locker
}

func NewReservationService(
	repo *repository.ReservationRepository,
	cacheClient cache.Client,
	locker lock.Locker,
) *ReservationService {
	return &ReservationService{
		repo:  repo,
		cache: cacheClient,
		lock:  locker,
	}
}

func (s *ReservationService) Reserve(
	ctx context.Context,
	resourceID string,
	start time.Time,
	end time.Time,
) error {
	lockKey := fmt.Sprintf(
		"lock:resource:%s:%d:%d",
		resourceID,
		start.Unix(),
		end.Unix(),
	)
	lockToken, err := s.lock.Acquire(ctx, lockKey, 30*time.Second)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	if lockToken == "" {

		metrics.ReservationsTotal.WithLabelValues("busy").Inc()

		return errors.New("resource busy, try again")
	}

	defer func() {
		if _, err := s.lock.Release(ctx, lockKey, lockToken); err != nil {
			fmt.Printf("failed to release lock: %v\n", err)
		}
	}()

	available, err := s.repo.CheckAvailability(ctx, resourceID, start, end)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	if !available {

		metrics.ReservationsTotal.WithLabelValues("conflict").Inc()

		return errors.New("resource not available")
	}

	err = s.repo.Create(ctx, resourceID, start, end)
	if err != nil {

		metrics.ReservationsTotal.WithLabelValues("error").Inc()

		return err
	}

	cacheKey := fmt.Sprintf(
		"availability:%s:%d:%d",
		resourceID,
		start.Unix(),
		end.Unix(),
	)
	_ = s.cache.Del(ctx, cacheKey)

	metrics.ReservationsTotal.WithLabelValues("success").Inc()

	return nil
}

func (s *ReservationService) GetReservations(
	ctx context.Context,
	resourceID string,
) ([]model.Reservation, error) {
	return s.repo.GetByResourceID(ctx, resourceID)
}
