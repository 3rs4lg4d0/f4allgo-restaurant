package service

import (
	"context"
	"errors"
	"fmt"

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

func (rs *DefaultRestaurantService) FindAll(ctx context.Context) ([]*domain.Restaurant, error) {
	//restaurants, err := rs.restaurantRepository.FindAll(ctx)
	//if err != nil {
	//	log.Error().Msg("an error occurred while fetching all the restaurant: " + err.Error())
	//	return nil, coreerrors.NewRepositoryError(err)
	//}

	return nil, nil
}

func (rs *DefaultRestaurantService) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		if err := rs.restaurantRepository.Save(ctx, restaurant); err != nil {
			log.Error().Msg("an error occurred while persisting the new created restaurant: " + err.Error())
			return coreerrors.NewRepositoryError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantCreated(restaurant)); err != nil {
			return fmt.Errorf("domainEventPublisher.Publish: %w", coreerrors.NewEventPublisherError(err))
		}

		return nil
	})
}

func (rs *DefaultRestaurantService) UpdateMenu(ctx context.Context, restaurantId uint64, menu *domain.Menu) error {
	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		restaurant, err := rs.findById(ctx, restaurantId)
		if err != nil {
			return err
		}

		if err := restaurant.UpdateMenu(menu); err != nil {
			return fmt.Errorf("restaurant.UpdateMenu: %w", err)
		}

		if err := rs.restaurantRepository.Update(ctx, restaurant); err != nil {
			return coreerrors.NewRepositoryError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantMenuUpdated(restaurantId, menu)); err != nil {
			return coreerrors.NewEventPublisherError(err)
		}

		return nil
	})
}

func (rs *DefaultRestaurantService) Delete(ctx context.Context, restaurantId uint64) error {
	return rs.trManager.Do(ctx, func(ctx context.Context) error {
		// We check the existence of the restaurant because we want to make sure we
		// don't send useless 'RestaurantDeleted' events.
		_, err := rs.findById(ctx, restaurantId)
		if err != nil {
			return err
		}

		if err := rs.restaurantRepository.Delete(ctx, restaurantId); err != nil {
			return coreerrors.NewRepositoryError(err)
		}

		if err := rs.domainEventPublisher.Publish(ctx, domain.NewRestaurantDeleted(restaurantId)); err != nil {
			return coreerrors.NewEventPublisherError(err)
		}

		return nil
	})
}

// findById tries to get a restaurant by its identifier using the repository. It
// also returns errors.NewRestaurantNotFoundError in case the restaurant doesn't
// exist.
func (rs *DefaultRestaurantService) findById(ctx context.Context, restaurantId uint64) (*domain.Restaurant, error) {
	restaurant, err := rs.restaurantRepository.FindById(ctx, restaurantId)
	if err != nil {
		return nil, coreerrors.NewRepositoryError(err)
	}

	if restaurant == nil {
		return nil, coreerrors.NewRestaurantNotFoundError(errors.New("restaurant not found"))
	}

	return restaurant, nil
}
