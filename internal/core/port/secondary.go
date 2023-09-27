package port

import (
	"context"
	"f4allgo-restaurant/internal/core/domain"
)

// RestaurantRepository manages persistent operations on restaurants dealing
// with an external storage system.
type RestaurantRepository interface {

	// FindAll restrieves all the registered restaurants with their menus.
	FindAll(ctx context.Context, offset uint32, limit uint8) ([]*domain.Restaurant, error)

	// FindById retrieves a particular restaurant by its identifier. This method
	// does not load the menu; just the restaurant info.
	FindById(ctx context.Context, restaurantId uint64) (*domain.Restaurant, error)

	// Save persists a restaurant in the external storage.
	Save(ctx context.Context, restaurant *domain.Restaurant) error

	// Update updates a restaurant and its relations. Call this operation when
	// you want to update a restaurant's menu.
	Update(ctx context.Context, restaurant *domain.Restaurant) error

	// Delete deletes a restaurant from the external storage.
	Delete(ctx context.Context, restaurantId uint64) error
}

// DomainEventPublisher publishes domain events to the outside world dealing with a
// message broker.
type DomainEventPublisher interface {

	// Publish publishes a domain event to the outside world.
	Publish(ctx context.Context, e domain.DomainEvent) error
}
