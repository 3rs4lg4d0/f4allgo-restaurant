package service

import (
	"context"

	"f4allgo-restaurant/internal/core/domain"
	"f4allgo-restaurant/internal/core/port"
	coreerrors "f4allgo-restaurant/internal/core/service/errors"

	"github.com/avito-tech/go-transaction-manager/trm"
	"github.com/rs/zerolog/log"
)

// Default implementation of the primary port RestaurantService.
type DefaultRestaurantService struct {
	restaurantRepository port.RestaurantRepository
	domainEventPublisher port.DomainEventPublisher
	trManager            trm.Manager
}

// Interface compliance verification.
var _ port.RestaurantService = (*DefaultRestaurantService)(nil)

// Provides an instance of DefaultRestaurantService.
func NewDefaultRestaurantService(restaurantRepository port.RestaurantRepository, domainEventPublisher port.DomainEventPublisher, trManager trm.Manager) *DefaultRestaurantService {
	return &DefaultRestaurantService{restaurantRepository: restaurantRepository, domainEventPublisher: domainEventPublisher, trManager: trManager}
}

func (rs *DefaultRestaurantService) FindAll(ctx context.Context, offset int, limit int) ([]*domain.Restaurant, int64, error) {
	restaurants, total, err := rs.restaurantRepository.FindAll(ctx, offset, limit)
	if err != nil {
		log.Error().Msg("an error occurred while fetching all the restaurant: " + err.Error())
		return nil, 0, coreerrors.NewRepositoryError(err)
	}

	return restaurants, total, nil
}

func (rs *DefaultRestaurantService) FindById(ctx context.Context, restaurantId int64) (*domain.Restaurant, error) {
	return rs.findById(ctx, restaurantId, true)
}

func (rs *DefaultRestaurantService) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		if err := rs.restaurantRepository.Save(ctx, restaurant); err != nil {
			log.Error().Msg("an error occurred while persisting the new created restaurant: " + err.Error())
			return coreerrors.NewRepositoryError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantCreated(restaurant)); err != nil {
			return coreerrors.NewEventPublisherError(err)
		}

		return nil
	})
}

func (rs *DefaultRestaurantService) UpdateMenu(ctx context.Context, restaurantId int64, menu *domain.Menu) error {
	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		restaurant, err := rs.findById(ctx, restaurantId, false)
		if err != nil {
			return err
		}

		if err := restaurant.UpdateMenu(menu); err != nil {
			return coreerrors.NewCoreError(err)
		}

		if _, err := rs.restaurantRepository.Update(ctx, restaurant); err != nil {
			return coreerrors.NewRepositoryError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantMenuUpdated(restaurantId, menu)); err != nil {
			return coreerrors.NewEventPublisherError(err)
		}

		return nil
	})
}

func (rs *DefaultRestaurantService) Delete(ctx context.Context, restaurantId int64) error {
	var rowsAffected int64
	var err error

	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		if rowsAffected, err = rs.restaurantRepository.Delete(ctx, restaurantId); err != nil {
			return coreerrors.NewRepositoryError(err)
		}

		if rowsAffected == 0 {
			return coreerrors.NewRestaurantNotFoundError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantDeleted(restaurantId)); err != nil {
			return coreerrors.NewEventPublisherError(err)
		}

		return nil
	})
}

// findById is a private function to find restaurants. It returns core service errors
// for any database access error and for restaurants not found in the database.
func (rs *DefaultRestaurantService) findById(ctx context.Context, restaurantId int64, fetchMenu bool) (*domain.Restaurant, error) {
	restaurant, err := rs.restaurantRepository.FindById(ctx, restaurantId, fetchMenu)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, coreerrors.NewRestaurantNotFoundError(err)
		} else {
			return nil, coreerrors.NewRepositoryError(err)
		}
	}

	return restaurant, nil
}
