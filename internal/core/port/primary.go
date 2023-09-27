package port

import (
	"context"

	"f4allgo-restaurant/internal/core/domain"
)

// RestaurantService exposes operations on restaurants. These operations are
// implemented in the service layer.
type RestaurantService interface {

	// FindAll gets the complete list of registered restaurants.
	FindAll(ctx context.Context, offset int, limit int) ([]*domain.Restaurant, error)

	// Create creates and persist a restaurant.
	Create(ctx context.Context, restaurant *domain.Restaurant) error

	// UpdateMenu updates a restaurant's menu enforcing some invariants.
	UpdateMenu(ctx context.Context, restaurantId uint64, menu *domain.Menu) error

	// Delete deletes a restaurant.
	Delete(ctx context.Context, restaurantId uint64) error
}
